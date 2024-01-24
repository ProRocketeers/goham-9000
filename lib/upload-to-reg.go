package lib

import (
	"fmt"

	"github.com/spf13/viper"
)

func Tst(uuid string, path string) string {
	// Here code co uploaduje do registru

	// ctx := context.Background()
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// if err != nil {
	// 	panic(err)
	// }

	// authConfig := types.AuthConfig{
	// 	Username:      "yourusername",
	// 	Password:      "yourpassword",
	// 	ServerAddress: "registry-url:port",
	// }
	// encodedJSON, err := json.Marshal(authConfig)
	// if err != nil {
	// 	panic(err)
	// }

	// authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	// options := types.ImagePushOptions{RegistryAuth: authStr}

	// responseBody, err := cli.ImagePush(ctx, "registry-url:port/yourimage:tag", options)
	// if err != nil {
	// 	panic(err)
	// }
	// defer responseBody.Close()

	// // Print the response to stdout
	// _, err = io.Copy(os.Stdout, responseBody)
	// if err != nil {
	// 	panic(err)
	// }

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	port := viper.Get("PORT")
	fmt.Println(port)

	return "with this uuid: " + uuid + " from this path: " + path
}
