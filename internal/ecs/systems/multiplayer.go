package systems

import (
	"encoding/gob"
	"game/internal/ecs"
	"game/internal/ecs/components"
	"game/internal/ecs/objects"
	"game/internal/network"
	"net"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

func init() {
	network.RegisterType(components.Transform2D{})
	network.RegisterType(components.Sprite{})
	network.RegisterType(components.Camera{})
	network.RegisterType(components.Input{})
	network.RegisterType(components.Network{})
	network.RegisterType(components.Player{})
	network.RegisterType(ecs.BaseComponent{})
}

var (
	multiplayerSystem = &MultiplayerSystem{}
	playerID          int
)

// -------------------- SPECIAL COMPONENTS -------------------------------------------//
type MultiplayerSystemRequiredComponents struct {
	Entities []*MultiplayerSystemComponents
}

type MultiplayerSystemComponents struct {
	Network *components.Network
}

// -------------------- Main Component -------------------------------------------//

type MultiplayerSystem struct {
	trackedEntities *MultiplayerSystemRequiredComponents

	serverConnection  *gob.Encoder
	clientConnections []*gob.Encoder
	isServer          bool
}

func NewMultiplayerSystem(isServer bool) *MultiplayerSystem {
	multiplayerSystem = &MultiplayerSystem{
		trackedEntities: &MultiplayerSystemRequiredComponents{},
		isServer:        isServer,
	}
	return multiplayerSystem
}

// -------------------- Custom Functionality -------------------------------------------//

func (pc *MultiplayerSystem) Update(dt float32) {
	// server
	// find structs from the network component, get their whole entity and find matches in the
	// entity with what should be uploaded every 60 ticks, every 4 ticks one should marshal all
	// known components to all connected players

	if pc.isServer {
		for _, entity := range pc.trackedEntities.Entities {
			id := entity.Network.GetComponentID()
			comps := ecs.GetActiveWorld().GetEntity(id)

			// acutally preform the replication
			for _, comp := range comps {
				// Handle Transform
				if transformComp, ok := comp.(*components.Transform2D); ok {
					for idx, connection := range pc.clientConnections {
						err := network.SendPacketTCP(connection, transformComp, "update comp")
						if err != nil {
							logrus.Error(err)
							if len(pc.clientConnections) > idx+1 && idx != 0 { // is it not first and not last i.e. middle
								pc.clientConnections = append(pc.clientConnections[:idx-1], pc.clientConnections[idx+1:]...)
							} else if idx != 0 { // is not first i.e. is it last
								pc.clientConnections = pc.clientConnections[:idx-1]
							} else if len(pc.clientConnections) > 1 { // is the anything other then it left i.e. its first
								pc.clientConnections = pc.clientConnections[1:]
							} else { // not in the middle, the last, or first position. no more positions left
								pc.clientConnections = []*gob.Encoder{}
							}
						}
					}
				}
			}
		}
	}

	// players
	// should get input componenets and other important data components and push them to the server
	// to replicate the changes

	if !pc.isServer {
		for _, entity := range pc.trackedEntities.Entities {
			id := entity.Network.GetComponentID()
			comps := ecs.GetActiveWorld().GetEntity(id)

			// see if it is a player TODO: SHOULD BE HANDLED BY NETWORK COMPONENT REPLICATION TYPE
			if comp, ok := comps[reflect.TypeOf(&components.Player{})]; ok {
				// Check if the comp is the player or not
				if comp.(*components.Player).PlayerID == playerID {

					// acutally preform the replication
					for _, comp := range comps {
						// Handle Input
						if inputComp, ok := comp.(*components.Input); ok {
							err := network.SendPacketTCP(pc.serverConnection, inputComp, "update comp")
							if err != nil {
								logrus.Error(err)
							}
						}
					}
				}
			}
		}
	}
}

func (pc *MultiplayerSystem) Draw(screen *ebiten.Image) {
	// Do nothing
}

func (pc *MultiplayerSystem) Initilizer() {

	// Start some network engine
	if pc.isServer {
		network.HandleTCPRequestFunc = pc.ServerHandleTCPRequest
		network.HandleUDPRequestFunc = pc.ServerHandleUDPRequest
		go network.StartListen()
	} else {
		network.HandleTCPResponseFunc = pc.ClientHandleTCPRequest

		conn, err := network.StartTCPConnection()
		if err != nil {
			logrus.Error(err)
		}

		pc.serverConnection = gob.NewEncoder(conn)
		network.SendPacketTCP(pc.serverConnection, "worldstate", "join request")
	}
}

func (pc *MultiplayerSystem) ClientHandleTCPRequest(enc *gob.Encoder, packet *network.Packet) {
	switch packet.Message.Type {

	case "update comp":
		logrus.Debugf("I got a update comp Packet: %v", packet)

		// make sure it actually is a component
		if comp, ok := packet.Message.Data.(ecs.Component); ok {
			id := comp.GetComponentID()
			entity := ecs.GetActiveWorld().GetEntity(id)

			comp = ecs.ToStructPtr(comp).(ecs.Component)
			compType := reflect.TypeOf(comp)

			currentComp := entity[compType]

			// Copy from comp to currentComp
			currentCompTransform, ok := currentComp.(*components.Transform2D)
			compTransform, ok2 := comp.(*components.Transform2D)
			if !ok2 || !ok {
				logrus.Info("breaking here to new comp")
				goto nextCase
			}

			currentCompTransform.WorldPosition = compTransform.WorldPosition
			currentCompTransform.WorldRotation = compTransform.WorldRotation
			currentCompTransform.WorldScale = compTransform.WorldScale

			break
		}
	nextCase:
		fallthrough

	case "new comp":
		logrus.Debugf("I got a worldstate Packet: %v", packet)
		logrus.Debugf("Data: %v", packet.Message.Data)
		data := packet.Message.Data
		switch data.(type) {
		case components.Player:
			playerComp := data.(components.Player)
			if playerComp.PlayerID == playerID { // TODO: Make AddComponent Copy value over if instance already exists
				newComponent, ok := data.(ecs.Component)
				if ok {
					val := reflect.ValueOf(newComponent)
					vp := reflect.New(val.Type())
					vp.Elem().Set(val)
					newComponent = vp.Interface().(ecs.Component)
					logrus.Debugf("Data After Handling: %v", newComponent)
					ecs.GetActiveWorld().AddComponent(newComponent)
				}
			} else {
				logrus.Infof("disregarding player %v", playerComp)
			}
		default:
			newComponent, ok := data.(ecs.Component)
			if ok {
				val := reflect.ValueOf(newComponent)
				vp := reflect.New(val.Type())
				vp.Elem().Set(val)
				newComponent = vp.Interface().(ecs.Component)
				logrus.Debugf("Data After Handling: %v", newComponent)
				ecs.GetActiveWorld().AddComponent(newComponent)
			}
		}

	case "id":
		logrus.Infof("ID Data: %v", packet.Message.Data)
		data := packet.Message.Data
		if id, ok := data.(int); ok {
			playerID = id
		} else {
			logrus.Warnf("id packet %v did not contain a valid int", data)
		}

	default:
		logrus.Warnf("type %v does not have a coded response for", packet.Message.Type)
	}
}

func (pc *MultiplayerSystem) ServerHandleTCPRequest(enc *gob.Encoder, packet *network.Packet) {
	logrus.Debug("Recieved TCP Packet %v", packet)

	switch packet.Message.Type {
	case "join request":
		// add new connections
		pc.clientConnections = append(pc.clientConnections, enc)

		playerID = playerID + 1

		// create their new player
		id := objects.NewPlayer(ecs.GetActiveWorld(), true, playerID)
		playerComp := ecs.GetActiveWorld().GetEntity(id)[reflect.TypeOf(&components.Player{})]
		logrus.Infof("current playerID: %v, respective Comp: %v", playerID, playerComp)

		err := network.SendPacketTCP(enc, playerID, "id")
		if err != nil {
			logrus.Error(err)
		}

		// send them all of your data
		for _, entity := range multiplayerSystem.trackedEntities.Entities {
			id := entity.Network.GetComponentID()
			comps := ecs.GetActiveWorld().GetEntity(id)
			for _, comp := range comps {
				switch comp.(type) {

				default: // just send by defualt
					logrus.Debugf("sending Data: %v", comp)

					err := network.SendPacketTCP(enc, comp, "new comp")
					if err != nil {
						logrus.Error(err)
					}

				}
			}
		}
	case "update comp":
		logrus.Debug("I got a worldstate Packet: %v", packet)

		// make sure it actually is a component
		if comp, ok := packet.Message.Data.(ecs.Component); ok {
			id := comp.GetComponentID()
			entity := ecs.GetActiveWorld().GetEntity(id)

			comp = ecs.ToStructPtr(comp).(ecs.Component)
			compType := reflect.TypeOf(comp)

			currentComp := entity[compType]

			// Copy from comp to currentComp
			currentCompInput := currentComp.(*components.Input)
			compInput := comp.(*components.Input)

			currentCompInput.Keys = compInput.Keys

		}

	default:
		logrus.Warnf("type %v does not have a coded response for", packet.Message.Type)
	}

}

func (pc *MultiplayerSystem) ServerHandleUDPRequest(conn net.Conn, packet *network.Packet) {
	logrus.Infof("Recieved UDP Packet %v", packet)
}

// -------------------- BoilerPlate Code -------------------------------------------//

func (pc *MultiplayerSystem) GetRequiredComponents() []reflect.Type {
	reqComponentsStruct := &MultiplayerSystemRequiredComponents{}

	v := reflect.ValueOf(reqComponentsStruct).Elem()

	returnTyppc := []reflect.Type{}
	for j := 0; j < v.NumField(); j++ {
		reqField := v.Field(j)
		switch reqField.Type().Kind() {
		case reflect.Slice:
			returnTyppc = append(returnTyppc, reqField.Type().Elem())
		case reflect.Ptr:
			returnTyppc = append(returnTyppc, reqField.Elem().Type())
		default:
			logrus.Error("no field match found")
		}
	}

	return returnTyppc
}

func (pc *MultiplayerSystem) AddEntity(comps map[reflect.Type]ecs.Component) {
	logrus.Trace("adding entity to Multiplayer System")

	for _, reqComp := range pc.GetRequiredComponents() {
		if ecs.SatisfySystemRequirements(comps, reqComp) {
			f := reflect.ValueOf(pc.trackedEntities).Elem()
			for j := 0; j < f.NumField(); j++ {
				reqField := f.Field(j)
				reqFieldType := reqField.Type().Elem()
				if reqFieldType == reqComp {
					newReqFieldEntry := reflect.New(reqFieldType.Elem())
					ecs.Fill(newReqFieldEntry, comps)

					reqFieldElem := reqField

					logrus.Debug("Setting Multiplayer System entity element")
					reqFieldElem.Set(reflect.Append(reqFieldElem, newReqFieldEntry))
				}
			}
		}
	}

}

func (pc *MultiplayerSystem) RemoveEntity(id ecs.ID) {
	// Called when a entity needs removed
}
