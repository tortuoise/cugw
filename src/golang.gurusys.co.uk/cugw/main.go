package main


import(

	"flag"
        "fmt"
        "github.com/grpc-ecosystem/grpc-gateway/runtime"
        "golang.org/x/net/context"
        "golang.gurusys.co.uk/go-framework/server"
	"golang.gurusys.co.uk/go-framework/utils"
        pbc "golang.gurusys.co.uk/apis/cugw"
        //pb "golang.gurusys.co.uk/apis/echoservice"
        "github.com/golang/protobuf/ptypes/timestamp"
        "github.com/golang/protobuf/proto"
        "github.com/golang/protobuf/ptypes"
        "github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"

	"golang.gurusys.co.uk/go-framework/client"
	"golang.gurusys.co.uk/go-framework/tokens"
)

var (
        _ pbc.AnythingServer = &Transparent{}
	port              = flag.Int("port", 10000, "The grpc server port")
)

func main() {

        t1 := &timestamp.Timestamp{Seconds:5, Nanos:10}
        serializedBytes, err := proto.Marshal(t1)
        if err != nil {
                fmt.Println("marshaling error: ", err)
        }
        fmt.Printf("Serialized timestamp: %d \n", len(serializedBytes))
        fmt.Printf("Message name: %s", proto.MessageName(t1))

        jsn := []byte("{\"Hello\":\"World\"}")

        s := &pbc.AnythingForYou{Anything: &any.Any {
                TypeUrl: "golang.gurusys.co.uk/" + proto.MessageName(t1),
                Value: serializedBytes }}
        jsonpb := &runtime.JSONPb{}
        err = jsonpb.Unmarshal(jsn, s)
        if err != nil {
                fmt.Println(err)
        }
        jsn,err = jsonpb.Marshal(s)
        fmt.Println(string(jsn))


	go func() {
		sd := server.NewServerDef()
		sd.Port = *port
		sd.Register = server.Register(
			func(server *grpc.Server) error {
				e := new(Transparent)
				pbc.RegisterAnythingServer(server, e)
				return nil
			},
		)
		err := server.ServerStartup(sd)
		utils.Bail("Unable to start server", err)
	}()

	con := client.Connect("cugw.Anything")
	anyClient := pbc.NewAnythingClient(con)
        ctx := tokens.ContextWithToken()
        //var a2 pbc.AnythingForYou
        var t2 timestamp.Timestamp
        resp, err := anyClient.Nothing(ctx, s)
        if err != nil {
                fmt.Println(err)
        }
        if err = ptypes.UnmarshalAny(resp.Anything, &t2); err != nil {
                fmt.Println("mismatch: ", err)
        }
        fmt.Printf("%#v \n", t2)

        //jsonpb.Marshal(

        ch := make(chan struct{})
        <-ch

}

type Transparent struct {

}


func (t *Transparent) Nothing(ctx context.Context, stuff *pbc.AnythingForYou) (*pbc.AnythingForYou, error) {

        return stuff, nil

}
