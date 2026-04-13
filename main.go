package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/govalues/decimal"
	"github.com/spf13/cobra"
)

var p, r float64
var ci, cd string

func main() {
	root := &cobra.Command{
		Use: "interest-calc",
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			end := shift(start, cd)
			step := shift(start, ci)

			totalHrs := end.Sub(start).Hours()
			stepHrs := step.Sub(start).Hours()
			periods := int(totalHrs / stepHrs)

			dailyRate := (r / 100) / 365
			periodicRate := dailyRate * (stepHrs / 24)

			pDec, _ := decimal.NewFromFloat64(p)
			rDec, _ := decimal.NewFromFloat64(periodicRate)
			base, _ := decimal.MustParse("1").Add(rDec)

			acc := decimal.MustParse("1")
			for range periods {
				acc, _ = acc.Mul(base)
			}

			res, _ := pDec.Mul(acc)
			final, _ := res.Float64()

			fmt.Printf("Projected: %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
			fmt.Printf("Periods: %d | Final: $%.2f | Interest: $%.2f\n", periods, final, final-p)
		},
	}

	root.Flags().Float64VarP(&p, "principal", "p", 0, "Principal amount")
	root.Flags().Float64VarP(&r, "rate", "r", 0, "Annual interest rate (as a percentage)")
	root.Flags().StringVarP(&ci, "interval", "i", "1d", "Interval of compounding (e.g., '1d', '1m', '1y')")
	root.Flags().StringVarP(&cd, "duration", "d", "1y", "Duration of the investment (e.g., '1y', '3m', '2d')")
	root.Execute()
}

func shift(t time.Time, s string) time.Time {
	v, _ := strconv.Atoi(s[:len(s)-1])
	switch s[len(s)-1:] {
	case "y":
		return t.AddDate(v, 0, 0)
	case "m":
		return t.AddDate(0, v, 0)
	case "w":
		return t.AddDate(0, 0, v*7)
	default:
		return t.AddDate(0, 0, v)
	}
}
