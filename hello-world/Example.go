package main

import (
	"encoding/json"
	"fmt"
	"github.com/global-soft-ba/decisionTable"
	"github.com/global-soft-ba/decisionTable/data"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"time"
)

type KnowledgeLib struct {
	Library *ast.KnowledgeLibrary
	Builder *builder.RuleBuilder
}

func CreateKnowledgeLibrary() *KnowledgeLib {
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	return &KnowledgeLib{knowledgeLibrary, ruleBuilder}
}

func (rb *KnowledgeLib) AddRule(rule string, knowledgeBase string) error {
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	bs := pkg.NewBytesResource([]byte(rule))
	err := rb.Builder.BuildRuleFromResource(knowledgeBase, "0.0.1", bs)
	if err != nil {
		return err
	}
	return nil
}

type Claim struct {
	TypeOfClaim        string
	ExpenditureOfClaim int
	TimeOfClaim        time.Time
}

type Employee struct {
	ResponsibleEmployee string
	FourEyesPrinciple   bool
	LastTime            time.Time
}

var first bool = true
var ruleEngine *engine.GruleEngine
var kb *ast.KnowledgeBase
var kl *KnowledgeLib

func RunRules(s string) string {
	ret := ""

	fmt.Println("RunRules Input = ", s)

	if first {
		table, _ := decisionTable.CreateDecisionTable().
			SetName("Determine Employee").
			SetDefinitionKey("determineEmployee").
			SetNotationStandard(data.GRULE).
			SetHitPolicy(data.Unique).
			AddInputField(data.TestField{Name: "Claim", Key: "TypeOfClaim", Typ: data.String}).
			AddInputField(data.TestField{Name: "Claim", Key: "ExpenditureOfClaim", Typ: data.Integer}).
			AddOutputField(data.TestField{Name: "Employee", Key: "ResponsibleEmployee", Typ: data.String}).
			AddOutputField(data.TestField{Name: "Employee", Key: "FourEyesPrinciple", Typ: data.Boolean}).
			AddRule("R1").
			AddInputEntry(`"Car Accident"`, data.SFEEL).
			AddInputEntry("<1000", data.SFEEL).
			AddOutputEntry(`"M??ller"`, data.SFEEL).
			AddOutputEntry("false", data.SFEEL).
			BuildRule().
			AddRule("R2").
			AddInputEntry(`"Car Accident"`, data.SFEEL).
			AddInputEntry("[1000..10000]", data.SFEEL).
			AddOutputEntry(`"Schulz"`, data.SFEEL).
			AddOutputEntry("false", data.SFEEL).
			BuildRule().
			AddRule("R3").
			AddInputEntry("-", data.SFEEL).
			AddInputEntry(">=10000", data.SFEEL).
			AddOutputEntry("-", data.SFEEL).
			AddOutputEntry("true", data.SFEEL).
			BuildRule().
			Build()

		// ConvertToGrlAst Table Into Grule Rules
		rules, err := table.Convert(string(data.GRULE))
		if err != nil {
			fmt.Print("Error:", err)
		}

		//Load Library and Insert rules
		fmt.Println("--------------GRL-RUlES------------------------")
		kl = CreateKnowledgeLibrary()
		result := rules.([]string)
		for _, rule := range result {
			fmt.Print(rule)
			addErr := kl.AddRule(rule, "#exampleBase")
			if addErr != nil {
				fmt.Print("Error:", addErr)
			}
		}
		// CreateEngine Instance
		ruleEngine = engine.NewGruleEngine()
		kb = kl.Library.NewKnowledgeBaseInstance("#exampleBase", "0.0.1")
		// Create Example Data
		//timeVal, err := time.Parse("2006-01-02T15:04:05", "2021-01-04T12:00:00")
		/*claim := Claim{
			TypeOfClaim:        "Car Accident",
			ExpenditureOfClaim: 100,
			TimeOfClaim:        timeVal,
		} */
		first = false
		/*
			f, err := os.OpenFile("myrules.grb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				err = kl.Library.StoreKnowledgeBaseToWriter(f, "myrules", "0.0.1")
				_ = f.Close()
				fmt.Print("OK writing rules file:")
			} else {
				fmt.Print("Error writing file:", err)
			}
		*/
	} else {
		/*
			ruleEngine = engine.NewGruleEngine()
			fmt.Print("OK created engine")
			kl := CreateKnowledgeLibrary()
			fmt.Print("OK created kl")
			f, err := os.Open("myrules.grb")
			if err != nil {
				fmt.Print("Error opening file:", err)
			} else {
				_, err := kl.Library.LoadKnowledgeBaseFromReader(f, true)
				if err != nil {
					fmt.Print("Error loading rules:", err)
				} else {
					fmt.Print("OK reading rules file:")
				}
				_ = f.Close()
			}
			_ = f.Close()
		*/
	}
	claim := Claim{}
	json.Unmarshal([]byte(s), &claim)
	employee := Employee{}

	// CreateEngine Instance
	//ruleEngine := engine.NewGruleEngine()
	//kb := kl.Library.NewKnowledgeBaseInstance("#exampleBase", "0.0.1")

	now := time.Now()
	// Load example
	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("Claim", &claim)
	err = dataCtx.Add("Employee", &employee)
	if err != nil {
		fmt.Println("Error:", err)
	}
	
	//Execution
	err = ruleEngine.Execute(dataCtx, kb)

	fmt.Println("--------------OutCome------------------------")
	fmt.Println("time elapse:", time.Since(now))

	fmt.Println("Input Claim =", claim)
	fmt.Println("Responsible Employee =", employee)
	out, err := json.Marshal(claim)
	if err != nil {
		panic(err)
	}
	out1, err := json.Marshal(employee)
	if err != nil {
		panic(err)
	}
	ret = string(out) + ":" + string(out1)
	return ret
}
