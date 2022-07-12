package main

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type (
	Order struct {
		Code    string  `json:"code"`
		Amount  float64 `json:"max_amount"`
		IsValid string  `json:"is_valid"`
	}
)

func (v *Order) ThisIsValid(sentence string) string {
	return fmt.Sprintf("Let say \"%s\"", sentence)
}

func main() {

	order := &Order{
		Code:   "AX456100",
		Amount: 700000,
	}

	// create data context
	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("VF", order)
	if err != nil {
		panic(err)
	}

	// create knowledge library
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	// create rule definition
	drls := `
	rule CheckVoucherValid "Check the voucher valid before July 20 2022 and amount less than 600K" salience 10 {
			when 
					VF.Code == "AX456100" &&  Now() < MakeTime(2022,10,20,0,0,0) && VF.Amount < 600000
			then
					VF.IsValid = VF.ThisIsValid("Congrat you got the discount!!");
					Retract("CheckVoucherValid");
	}
	`

	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	bs := pkg.NewBytesResource([]byte(drls))
	err = ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	// create Knowledge Base instance
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	// create Rule Engine instance
	engine := engine.NewGruleEngine()
	err = engine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	fmt.Println(order.IsValid)

}
