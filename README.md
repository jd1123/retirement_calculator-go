Advanced Retirement Calculator
=============================

Why make this?
-------

Traditional wisdom says that if you max your 401k and manage to get 7% returns per year you should be okay in retirement.
The real world is not like this. The stock market has ups and downs, and volatility can ruin performance.
I wanted to see how my portfolio would act with some more real world scenarios. I cannot find something online that treats this with any meaningful rigor.

A model is only as good as it's assumptions. Therefore, any way we can use more reaonsable assumptions in your retirement planning, the more informed you are.
I do not like all the retirement portfolio planners I see. They all assume constant returns over a period of time. You can set this return, but it does not
give a reasonable path for your nest egg. Furthermore, there is no "what-if" planning, i.e. how much more confidence does saving an extra $1,000 per year give 
me that I can retire comfortably?

What does it do differently?
------
The stock market is said to follow a [geometric brownian motion](https://en.wikipedia.org/wiki/Geometric_Brownian_motion) with a positive return and some [volatility](http://www.commonwealth.com/RepSiteContent/stock_volatility.htm). As a result, returns will be random.
We can simulate stock market data by using random data from a [N(mu,sigma) distribution](https://en.wikipedia.org/wiki/Log-normal_distribution). 
It is widely known that the stock market is [not lognormally distributed](https://en.wikipedia.org/wiki/Volatility_smile), but this is a much better approximation than assuming 7% smooth returns in perpetuity.
This is important, because your life and retirement plan are [path dependent](https://en.wikipedia.org/wiki/Path_dependence). There could be situations where a higher average return is actually worse for your 
savings plan.
 
How do I use this?
------
Right now its a work in progress. I am going to put this into a ~~django~~ Go webapp as soon as I find the time and figure out all the details. Currently, 
I am working out bugs in the logic.

The idea will be to [simulate](https://en.wikipedia.org/wiki/Monte_Carlo_method) ~~10000~~ 100000 paths of the stock market over any number of years and look at the distribution of your retirement portfolio.
You will then be able to plan with a certain confidence level.
Perhaps you're in good shape, but with what probability can you assume this is true?

One day, when I've worked on it enough, I will post it on a website and you will be able to use it remotely. Now, you can build it and run it locally:

```
mkdir -p <working_directory>
git clone https://github.com/jd1123/retirement_calculator-go
cd retirement_calculator-go
make
./retirement_calculator-go
```

Now point your browser to [127.0.0.1:8081](http://127.0.0.1:8081) and have fun.

Does this make any assumtions?
-----
Yes, it does. 

1. Savings for a particular year make no returns. This is not practically true.
2. Your retirement income does not adjust for inflation (I need to check this, it's been a long time since I wrote this).
3. Your taxable accounts are taxed at the long term capital gains tax each year. This may or may not be the case for you, but either way, this assumption was made to simulate tax liabilities stemming from returns generation.
4. Your non taxable accounts are taxed at 30% (check this number) upon withdrawal. Again, this is YOU specific, but in the future I'd like to make this an editable field.
5. The stock market follows a lognormal distrubution with Mu=7% and Sigma=15%. You can edit these by changing the type of portfolio: High Risk, Medium Risk or Low Risk.

Why re-write this in Go?
-----
Two reasons. One, I was running into memory issues with Python. Two, I want to learn Go.

Who are you?
------
I am johnnydiabetic. I am not a financial advisor, but I have worked in finance for more than a decade, first in risk management, then as an options trader. Please note, this does not in any way make me qualified to give financial advice. This is intended to be used as a tool to verify the impact using the stock market as a savings vehicle.
