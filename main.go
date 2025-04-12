package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const ProjectURL = "https://github.com/ivansouzamf/MathBot"
const DevUserID = "415179296823312394"

const Prefix = ";; "

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No token provided")
	}

	log.Println("Hello 👋. Starting Math Bot")
	
	Token := os.Args[1]
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("Couldn't create discord session")
	}

	session.AddHandler(messageCallback)

	// we only care about server messages
	session.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentMessageContent

	err = session.Open()
	if err != nil {
		log.Fatal("Couldn't connect to the discord API")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// cleanup
	session.Close()
	log.Println("Math Bot exited sucessfully")
}

func messageCallback(session *discordgo.Session, message *discordgo.MessageCreate) {
	// ignore messages of itself && non prefixed
	if message.Author.ID == session.State.User.ID || !strings.HasPrefix(message.Content, Prefix) {
		return
	}
	
	log.Println("Message received:", message.Content)
	
	if strings.Contains(message.Content, "about") {
		about := "I'm a bot capable of solving mathematical expressions\n"
		about += "Developed by <@" + DevUserID + "> in " + ProjectURL
		
		session.ChannelMessageSend(message.ChannelID, about)
	} else {
		exp := strings.ReplaceAll(message.Content, Prefix, "")
		res, err := EvaluateMathExp(exp)
		
		var answer string
		if err != nil {
			answer = "Invalid math expression"
		} else {
			answer = fmt.Sprintf("%s = %f", exp, res)
		}

		session.ChannelMessageSend(message.ChannelID, answer)
	}
}
