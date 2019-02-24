package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/embed"
)

func commandRoll(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	var dices []int

	for _, a := range args {
		tmp := strings.Split(a, "d")
		var err error
		var number int
		var max int
		if len(tmp) == 1 {
			number = 1
			max, err = strconv.Atoi(tmp[0])
		} else {
			number, err = strconv.Atoi(tmp[0])
			max, err = strconv.Atoi(tmp[1])
		}
		if err != nil {
			return nil, "Mauvaise requête"
		}
		dices = append(dices, createDices(number, max)...)
	}

	if len(dices) > 50 {
		return nil, "Wesh, calme-toi et demande moins de 50 dés."
	}

	numbers := readLines(b.RandomNumbers)
	if len(numbers) < len(dices) {
		str := callRandomOrg(b.BotConfig.RandomToken, b.RandomNumbers)
		if len(str) > 0 {
			return nil, str
		}
		numbers = readLines(b.RandomNumbers)
	}

	var values []int
	for i := 0; i < len(dices); i++ {
		d := numbers[i]*float64(dices[i]) + 1
		values = append(values, int(math.Round(d)))
	}

	writeNumbers(numbers[len(dices):], b.RandomNumbers)
	return formatDices(dices, values, env.Message.Author.Mention())
}

func formatDices(dices []int, values []int, user string) (*discordgo.MessageEmbed, string) {

	fields := []*discordgo.MessageEmbedField{}
	res := make(map[int][]int)
	for i := 0; i < len(values); i++ {
		res[dices[i]] = append(res[dices[i]], values[i])
	}

	for key, value := range res {
		fields = append(fields, embed.EmbedField(strconv.Itoa(key), fmt.Sprint(value), true))
	}

	return &discordgo.MessageEmbed{
		Title:       "Voici votre tirage",
		Description: "Tirage demandé par " + user,
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Randomisation fourni par random.org",
		},
	}, ""
}

func createDices(number int, max int) []int {
	var res []int
	for i := 0; i < number; i++ {
		res = append(res, max)
	}
	return res
}

func readLines(path string) []float64 {
	file, _ := os.Open(path)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	var numbers []float64
	for _, elem := range lines {
		i, err := strconv.ParseFloat(elem, 64)
		if err == nil {
			numbers = append(numbers, i)
		}
	}
	return numbers
}

func callRandomOrg(token string, db string) string {
	url := "https://api.random.org/json-rpc/1/invoke"
	params := Params{
		APIKey:        token,
		N:             5,
		DecimalPlaces: 5,
	}

	request := Request{
		Jsonrpc: "2.0",
		Method:  "generateDecimalFractions",
		Params:  params,
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
	fmt.Println("Data:", data)
	fmt.Println("Failed:", data.Failed)
	fmt.Println("Message:", data.Failed.Message)
	if data.Result.BitsUsed == 0 {
		return "Error: " + data.Failed.Message
	}

	writeNumbers(data.Result.Random.Data, db)

	file, _ := os.Create(db)

	for _, nb := range data.Result.Random.Data {
		fmt.Fprintln(file, nb)
	}
	file.Close()
	return ""
}

func writeNumbers(numbers []float64, path string) {
	file, _ := os.Create(path)

	for _, nb := range numbers {
		fmt.Fprintln(file, nb)
	}
	file.Close()
}

type Request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	APIKey        string `json:"apiKey"`
	N             int    `json:"n"`
	DecimalPlaces int    `json:"decimalPlaces"`
}

type RandomOrg struct {
	Jsonrpc string `json:"jsonrpc"`
	Failed  Failed `json:"error"`
	Result  Result `json:"result"`
}

type Result struct {
	Random        Random `json:"random"`
	BitsUsed      int    `json:"bitsused"`
	BitsLeft      int    `json:"bitsleft"`
	RequestsLeft  int    `json:"requestsleft"`
	AdvisoryDelay int    `json:"advisorydelay"`
}

type Random struct {
	Data           []float64 `json:"data"`
	CompletionTime string    `json:"completiontime"`
}

type Failed struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
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
