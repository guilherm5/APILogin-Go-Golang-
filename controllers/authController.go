package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/guilherm5/GoAndReact/database"
	"github.com/guilherm5/GoAndReact/models"
	"golang.org/x/crypto/bcrypt"
)

// secret key para jwt
const SecretKey = "secret"

// criar usuario
func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err, " error")
	}

	//encripitando senha do usuario

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		log.Println("Erro ao criptografar senha ", err)
	}

	//criando usuario
	user := models.UserAuth{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}
	database.DB.Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.JSON(http.ErrAbortHandler.Error())
		log.Println(err, " error")
	}

	//capturando usuario a partir do email
	var user models.UserAuth

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)

		return c.JSON(fiber.Map{
			"message": "usuario nao encontrato",
		})
	}

	//comparando senha
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "passowrd incorreto",
		})

	}

	//gerando jwt caso a senha esteja correta

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		log.Println("erro ao realizar login ", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "erro ao realizar login",
		})

	}

	//retornando mensagem de sucesso e o jwt indo para o cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})

}

//verificando qual usuario esta logado

func User(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		log.Println("erro ao realizar login ", err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "erro ao autenticar",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var user models.UserAuth

	database.DB.Where("id = ?", claims.Issuer).First(&user)
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "sucesso ao deslogar",
	})
}
