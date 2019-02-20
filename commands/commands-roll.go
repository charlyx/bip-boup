package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/bot"
)

func commandRoll(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	var dices []int
	for _, a := range args {
		tmp := strings.Split(a, "d")
		var number string
		var max string
		if len(tmp) == 1 {
			number = "1"
			max = tmp[0]
		} else {
			number = tmp[0]
			max = tmp[1]
		}
		dices = append(dices, callRandomOrg(number, max, env.Message.Author.ID, b.BotConfig.RandomToken)...)
	}

	res := env.Message.Author.Mention()
	for _, d := range dices {
		res += " " + strconv.Itoa(d)
	}
	return nil, res
}

func callRandomOrg(n string, max string, id string, token string) []int {
	fmt.Println("send:", n, "size", max)
	url := "https://api.random.org/json-rpc/1/invoke"
	params := Params{
		APIKey: token,
		N:      n,
		Min:    "1",
		Max:    max,
	}

	request := Request{
		Jsonrpc: "2.0",
		Method:  "generateIntegers",
		Params:  params,
		ID:      id,
	}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(request)

	req, err := http.NewRequest("POST", url, buffer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	randomOrg, _ := ioutil.ReadAll(resp.Body)

	var data RandomOrg
	json.Unmarshal(randomOrg, &data)
	fmt.Println("Receive:", data.Result.RequestsLeft)
	return data.Result.Random.Data
}

type Request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      string `json:"id"`
}

type Params struct {
	APIKey string `json:"apiKey"`
	N      string `json:"n"`
	Min    string `json:"min"`
	Max    string `json:"max"`
}

type RandomOrg struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
	ID      string `json:"id"`
}

type Result struct {
	Random        Random `json:"random"`
	BitsUsed      int    `json:"bitsused"`
	BitsLeft      int    `json:"bitsleft"`
	RequestsLeft  int    `json:"requestsleft"`
	AdvisoryDelay int    `json:"advisorydelay"`
}

type Random struct {
	Data           []int  `json:"data"`
	CompletionTime string `json:"completiontime"`
}

func init() {
	commands["roll"] = &bot.Command{
		Function: commandRoll,
		HelpText: "Lance pour vous des dés",
		Arguments: []bot.CommandArgument{
			{Name: "requête", Description: "Une requête sous la forme xdy, x étant le nombre de dés (vide accepté) et y le type de dés", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}
}
