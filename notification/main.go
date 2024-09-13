package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "notification/config"
)

func main() {
    app := fiber.New()

    // RabbitMQ bağlantısını başlat
    conn := config.ConnectRabbitmq()
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    // RabbitMQ'ya mesaj gönderen endpoint
    app.Post("/send", func(c *fiber.Ctx) error {
        message := c.FormValue("message")
        config.PublishMessage(ch, "myqueue", message)
        return c.SendString("Message sent: " + message)
    })

    // RabbitMQ'dan mesajları tüketen işlem
    go config.ConsumeMessages(ch, "myqueue")

    log.Fatal(app.Listen(":9090"))
}