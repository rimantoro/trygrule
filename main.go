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
		Amount     int64
		MerchantID string
	}

	// Rule attributes - Will be compared with Fact attributes within a rule
	RuleAttr struct {
		Type               string
		MinimumOrderAmount int64
		DiscAmount         int64
		IsActive           bool
		MerchantID         string
		Priority           int64
	}

	// This one is the Fact attributes
	OrderFact struct {
		Order    Order
		RuleAttr RuleAttr
		// Ini utk nampung result
		Eligible    bool
		AppliedDisc int64
	}
)

/**
Basic need dari OrderFact. Kita mau tahu apakah eligible, kalau iya, amountnya berapa
**/

func (of *OrderFact) IsEligible() bool {
	return of.Eligible
}

func (of *OrderFact) GetDiscAmount() int64 {
	return of.AppliedDisc
}

func main() {

	var (
		err error
	)

	// Semua fakta yang ada ditampung di struct ini. Order yg mau dicek + Parameter terkait discount (RuleAttr)
	of01 := &OrderFact{
		Order: Order{
			Amount:     150000,
			MerchantID: "MERCH002",
		},
		RuleAttr: RuleAttr{
			MinimumOrderAmount: 150000,
			DiscAmount:         25000,
			IsActive:           true,
			MerchantID:         "MERCH001",
		},
	}

	dataCtx := ast.NewDataContext()

	err = dataCtx.Add("OF", of01)
	if err != nil {
		panic(err)
	}

	/**
	  SETUP KNOWLEDGE LIBRARY
	  **/

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	// Rule attributes and Fact attributes are validate here using Grule script
	ruleQry := `
	rule CheckEligible "Check if eligible" salience 10 {
	when 
			OF.Order.Amount >= OF.RuleAttr.MinimumOrderAmount && 
            OF.Order.MerchantID == OF.RuleAttr.MerchantID && 
            OF.RuleAttr.IsActive == true
	then
			OF.Eligible=true;
			OF.AppliedDisc=OF.RuleAttr.DiscAmount;
			Retract("CheckEligible");
	}

	rule CheckNotEligible "Check if not eligible" salience 9 {
	when 
			OF.Order.Amount < OF.RuleAttr.MinimumOrderAmount ||
			OF.Order.MerchantID != OF.RuleAttr.MerchantID ||
			OF.RuleAttr.IsActive != true
	then
			OF.Eligible=false;
			OF.AppliedDisc=0;
			Retract("CheckNotEligible");
	}
  `

	bs := pkg.NewBytesResource([]byte(ruleQry))
	err = ruleBuilder.BuildRuleFromResource("FixShippingDiscRules", "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	/**
	  EXECUTE RULE ENGINE
	  **/

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("FixShippingDiscRules", "0.0.1")
	engine := engine.NewGruleEngine()
	err = engine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	fmt.Println(of01.IsEligible(), of01.GetDiscAmount())

}
