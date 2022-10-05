package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	ce "github.com/finacore/commons-errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Test struct {
	Name       string
	Email      string
	Quantidade int
	Valor      float32
}

// ConnectDB returns a tuple (connection, error) that represents the
func CreateConnection(mongoURI string, timeout int) (*mongo.Client, error) {
	var err error
	var client *mongo.Client

	//set options and create a new mongodb client. in case of unsuccess client
	//creation, this action finalize the program.
	var options *options.ClientOptions = options.Client().ApplyURI(mongoURI)
	if client, err = mongo.NewClient(options); err != nil {
		log.Print(fmt.Errorf("create client error: %w", err))
		return nil, err
	}

	//crete context and established connection. The connection established is
	//mandatory to the program execution. If it not possible the program must be
	//end with error status
	connectionTimeout := time.Duration(timeout) * time.Second
	ctx, _ := context.WithTimeout(context.Background(), connectionTimeout)
	if err = client.Connect(ctx); err != nil {
		log.Print(fmt.Errorf("make connection error: %w", err))
		return nil, err
	}

	//this is just a check to assert that the connection is working proper. If's
	//not, an error must be raised and the program will end with error status
	if err = client.Ping(ctx, nil); err != nil {
		log.Print(fmt.Errorf("assert connection error: %w", err))
		return nil, err
	}

	return client, nil
}

type ErrorMap struct {
	Reason       string
	OperatorName string
	PropertyName string
}

func (em *ErrorMap) String() string {
	return fmt.Sprintf(
		"%s for %s value",
		em.Reason,
		em.OperatorName,
	)
}

func (em *ErrorMap) IsComplete() bool {
	if em.OperatorName == "" {
		return false
	}

	if em.PropertyName == "" {
		return false
	}

	if em.Reason == "" {
		return false
	}

	return true
}

func (em *ErrorMap) Error() *ce.ValidationError {
	return ce.CreateValidationError(
		em.PropertyName,
		em.String(),
	)
}

// func E121ParseError(data map[string]interface{}, errorMap *ErrorMap) {
// 	if data["operatorName"] != nil {
// 		errorMap.OperatorName = data["operatorName"].(string)

// 		errorMap.OperatorName = data["operatorName"].(string)
// 	}

// 	if data["reason"] != nil {
// 		errorMap.Reason = data["reason"].(string)
// 	}

// 	if data["propertyName"] != nil {
// 		errorMap.PropertyName = data["propertyName"].(string)
// 	}

// 	for _, v := range data {
// 		switch v.(type) {
// 		case string:
// 			continue
// 		case []interface{}:
// 			vet := v.([]interface{})
// 			for i := range vet {
// 				E121ParseError(vet[i].(map[string]interface{}), errorMap)
// 			}
// 		case interface{}:
// 			E121ParseError(v.(map[string]interface{}), errorMap)
// 		}
// 	}
// }

// func e121ParseMap(mapa map[string]interface{}, errs *[]*ce.ValidationError) {
// 	for
// }

// func E121ParserArray(data []interface{}, errs *[]*ce.ValidationError) {
// 	for i := range data {
// 		switch data[i].(type) {
// 		case string:
// 			continue
// 		case []interface{}:
// 			E121ParserArray(data[i].([]interface{}), errs)
// 		case map[string]interface{}:

// 		}
// 		// E121Parser(vet[value].(map[string]interface{}), errs)
// 	}
// }

// func E121Parser(data map[string]interface{}, errs *[]*ce.ValidationError) {
// 	propertyName := data["propertyName"]
// 	description := data["description"]

// 	if propertyName != nil {
// 		if description != nil {
// 			*errs = append(
// 				*errs,
// 				ce.CreateValidationError(propertyName.(string), description.(string)),
// 			)
// 			return
// 		} else {
// 			fmt.Printf("%v\n", data)
// 			var e ErrorMap
// 			E121ParseError(data, &e)
// 			*errs = append(
// 				*errs,
// 				e.Error(),
// 			)
// 		}
// 	}

// 	for _, v := range data {
// 		switch v.(type) {
// 		//this is here because string is an interface. if this 'case' expression was
// 		//removed, the next one 'interface{}' wil be resolved and the program will try
// 		//to cast the string as map[string]interface{}.
// 		case string:
// 			continue
// 		case []interface{}:
// 			vet := v.([]interface{})
// 			for value := range vet {
// 				E121Parser(vet[value].(map[string]interface{}), errs)
// 			}
// 		case interface{}:
// 			mapa := v.(map[string]interface{})

// 			E121Parser(mapa, errs)

// 		}
// 	}
// }

type E121 struct {
	Errors []*ce.ValidationError
	Runner []interface{}
}

func NewE121(root interface{}) *E121 {
	e := &E121{}

	e.Runner = append(e.Runner, root)

	return e
}

func (e *E121) Run() {

	for len(e.Runner) > 0 {
		switch e.Runner[0].(type) {
		case string:
			//to nothing
		case []interface{}:
			vet := e.Runner[0].([]interface{})
			for j := range vet {
				e.Runner = append(e.Runner, vet[j])
			}
		case interface{}:
			data := e.Runner[0].(map[string]interface{})
			propertyName := data["propertyName"]
			description := data["description"]

			if propertyName != nil {
				if description != nil {
					e.Errors = append(e.Errors, ce.CreateValidationError(propertyName.(string), description.(string)))
				} else {
					fmt.Println(" #####  Nao achou o description")
					details := data["details"].([]interface{})
					var reason string

					for index := range details {
						switch details[index].(type) {
						case map[string]interface{}:
							item := details[index].(map[string]interface{})
							reason = item["reason"].(string)
						}
					}

					e.Errors = append(e.Errors, ce.CreateValidationError(propertyName.(string), reason))
				}
				// } else if reason != nil {
				// 	fmt.Printf("\t\n @reason >> %s", reason)
			} else {
				for _, v := range data {
					fmt.Printf("Add element >>> %v\n", v)
					e.Runner = append(e.Runner, v)
				}
			}
		}
		e.Runner = e.Runner[1:]
		fmt.Printf("\t### len >> %d\n", len(e.Runner))
	}
}

// func e121Interface(data map[string]interface{}, errs *[]*ce.ValidationError) {
// 	for k, v := range data {

// 	}
// }

func e121(err mongo.WriteException) []*ce.ValidationError {
	str := err.Raw.String()

	var dto map[string]interface{}
	json.Unmarshal([]byte(str), &dto)

	fmt.Printf("%v\n", dto)

	r := NewE121(dto)
	r.Run()

	print(len(r.Errors))
	// E121Parser(dto, &errs)

	return r.Errors

}

func main() {
	client, err := CreateConnection("mongodb://localhost:27017", 10)
	if err != nil {
		panic("connection error")
	}

	collection := client.Database("test").Collection("test")

	model := &Test{
		Name:       "marcos j",
		Email:      "joao@silva.com",
		Quantidade: 6,
		Valor:      0,
	}

	if _, err := collection.InsertOne(context.TODO(), model); err != nil {
		var myError mongo.WriteException
		errors.As(err, &myError)

		ve := e121(myError)

		for i := range ve {
			log.Println(ve[i].Error())
		}
	}
}
