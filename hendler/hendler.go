package hendler

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type Finance struct {
	ID    int    `json:"id"`
	Sum   int32  `json:"sum"`
	Types string `json:"type"`
}

var DB *sql.DB

func init() {

	err := godotenv.Load("/home/kirill/Gohka/finance/.env")
	if err != nil {
		log.Println("error load file")
	} else {
		fmt.Println("Successfully loadeed env")
	}

	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	host := os.Getenv("host")
	port := os.Getenv("port")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Усшеное подлючение")
}

// @Summary		Get all finances
// @Description	Get all finances
// @Tags			finances
// @Accept			json
// @Produce			json
// @Success		200	{object}	hendler.Finance
// @Router			/all [get]
func GetAll(c *gin.Context) {
	res, err := DB.Prepare("SELECT * FROM financeprod")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}
	defer res.Close()

	rows, err := res.Query()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}
	defer rows.Close()

	financ := []Finance{}

	for rows.Next() {
		fin := Finance{}
		if err := rows.Scan(&fin.ID, &fin.Sum, &fin.Types); err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
			return
		}
		financ = append(financ, fin)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, financ)
}

// @Summary		Add finance
// @Description	Add a new finance record
// @Tags			finances
// @Accept			json
// @Produce		json
// @Param 			finance body hendler.Finance true "Finance data to add"
// @Success		200	{object}	hendler.Finance
// @Router			/addfin [post]
func Postfinc(c *gin.Context) {
	var json struct {
		Summ  int32  `json:"sum"`
		Types string `json:"type"`
	}

	// Пробуем считать JSON из тела запроса
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println("Ошибка при парсинге данных:", err)
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	// Логируем значения данных перед вставкой в базу
	log.Printf("Полученные данные: sum=%d, type=%s", json.Summ, json.Types)

	// Подготовка запроса с использованием плейсхолдеров
	stmt, err := DB.Prepare("INSERT INTO financeprod (sum, type) VALUES ($1, $2)") // PostgreSQL использует $1, $2
	if err != nil {
		log.Println("Ошибка подготовки запроса:", err)
		c.JSON(500, gin.H{"error": "Ошибка при подготовке запроса"})
		return
	}

	// Выполнение запроса с передачей значений
	_, err = stmt.Exec(json.Summ, json.Types)
	if err != nil {
		log.Println("Ошибка выполнения запроса:", err)
		c.JSON(500, gin.H{"error": "Ошибка при выполнении запроса"})
		return
	}

	// Успешный ответ
	c.JSON(200, gin.H{"message": "Данные успешно добавлены"})
}

// @Summary		Get finance by ID
// @Description	Get finance by ID
// @Tags			finances
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Finance ID"
// @Success		200	{object}	hendler.Finance
// @Router			/finance/{id} [get]
func IdSearch(c *gin.Context) {
	var rep Finance
	uid := c.Param("id")

	log.Println("ID ==", uid)

	resp, err := DB.Prepare("SELECT  * FROM financeprod Where financeprod.id = $1")
	if err != nil {
		log.Println("ошибка выполнения запроса", err)
		c.JSON(500, gin.H{"error": "ошибка выполения запроса"})
	}

	defer resp.Close()

	err = resp.QueryRow(uid).Scan(&rep.ID, &rep.Sum, &rep.Types)
	if err != nil {
		log.Println("что-то не так", err)
		c.JSON(500, gin.H{"error": "данные не получены"})
	}

	c.JSON(200, rep)

}

type FinanceUpdateRequest struct {
	Sum   int32  `json:"sum" example:"200"`
	Types string `json:"type" example:"expense"`
}

type FinanceResponse struct {
	Message string `json:"message" example:"Данные успешно обновлены"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Некорректные данные"`
}

// @Summary		Update finance
// @Description	Update finance by ID
// @Tags			finances
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Finance ID"
// @Param			finance body FinanceUpdateRequest true "Finance data to update"
// @Success		200	{object}	FinanceResponse
// @Failure		400	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router			/updatefin/{id} [put]
func Updatefin(c *gin.Context) {
	var updfin FinanceUpdateRequest
	par := c.Param("id")

	log.Println("id с параметров", par)

	if err := c.ShouldBindJSON(&updfin); err != nil {
		log.Println("Не удалось прочитать json", err)
		c.JSON(400, ErrorResponse{Error: "Некорректные данные"})
		return
	}

	request, err := DB.Prepare("UPDATE financeprod SET sum = $1, type = $2 WHERE id = $3")
	if err != nil {
		log.Println("Ошибка формирования запроса", err)
		c.JSON(500, ErrorResponse{Error: "Ошибка при подготовке запроса"})
		return
	}

	_, err = request.Exec(updfin.Sum, updfin.Types, par)
	if err != nil {
		log.Println("Ошибка обновления данных", err)
		c.JSON(500, ErrorResponse{Error: "Ошибка при обновлении данных"})
		return
	}

	c.JSON(200, FinanceResponse{Message: "Данные успешно обновлены"})
}

// @Summary		Remove finance
// @Description	Remove finance by ID
// @Tags			finances
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Finance ID"
// @Success		200	{object}	hendler.Finance
// @Router			/removefin/{id} [delete]
func RemoveRecord(c *gin.Context) {
	resp := c.Param("id")

	log.Println("Id удаляемой записи", resp)

	stwe, err := DB.Prepare("DELETE FROM financeprod Where financeprod.id = $1")
	if err != nil {
		log.Println("Запрос не сформирован", err)
		c.JSON(500, gin.H{"error": "не удалось сформировать запрос"})
	}

	_, err = stwe.Exec(resp)
	if err != nil {
		log.Println("не у далось выполнить запрос", err)
		c.JSON(500, gin.H{"error": "заплос не выполнен"})
	}

	c.JSON(200, gin.H{"message": "запись удалена"})

}
