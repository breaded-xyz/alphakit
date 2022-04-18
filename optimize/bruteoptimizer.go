package optimize

import (
	"context"
)

type BruteOptimizer struct {
	testCases []CartesianProduct
}

func NewBruteOptimizer() *BruteOptimizer {
	return nil
}

func (o *BruteOptimizer) Configure(map[string]any) error {
	return nil
}

func (o *BruteOptimizer) Prepare(params map[string]any) (int, error) {
	o.testCases = CartesianBuilder(params)
	return len(o.testCases), nil
}

func (o *BruteOptimizer) Start(ctx context.Context) (chan<- StepResult, error) {

	return nil, nil
}

/*wp := workerpool.New(16)
var mu sync.Mutex
for i := range testCases {
	i := i
	wp.Submit(func() {
		tCase := testCases[i]
		dealer := backtest.NewDealer()
		if err := dealer.Configure(tCase); err != nil {
			if errors.Is(err, trader.ErrInvalidConfig) {
				return
			}
			panic(err)
		}

		bot := _typeRegistry[config["bot"].(string)].(botMakerFunc)()
		bot.SetDealer(dealer)
		if err := bot.Configure(tCase); err != nil {
			if errors.Is(err, trader.ErrInvalidConfig) {
				return
			}
			panic(err)
		}

		result, err := execBacktest(bot, dealer, prices)
		if err != nil {
			panic(err)
		}
		result.Description = fmt.Sprintf("%+v", tCase)

		mu.Lock()
		results = append(results, result)
		mu.Unlock()

		bar.Add(1)
	})
}

wp.StopWait()*/
