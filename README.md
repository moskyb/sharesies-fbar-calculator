# Disclosure!

I'm some third party rando. I used to work at Sharesies, but i don't anymore. This software comes as is, with no warranty, etc, and i'm not liable for anything that happens to you or your money. I'm no longer affiliated with Sharesies in any way, and i don't speak for them. I'm just some guy who made this software for to help with (sigh) filing my FBARs.

This program and the information it produces are not financial advice. You should consult a professional before making any financial decisions, including (especially!) filing taxation information with the US Government.

This program also may contravene the Sharesies Terms of Service. In using it, you recognise that Sharesies may decide to terminate your account for breach of Terms of Service, and that that's your fault, not mine.

# Sharesies FBAR Calculator

The US Department of the Treasury requires US citizens and residents to report their foreign financial accounts to the IRS. This includes bank accounts, brokerage accounts, mutual funds, etc. The FBAR is a form that must be filed electronically with the IRS. The FBAR is due on June 30th of each year, and covers the previous calendar year.

One such applicable financial account is the New Zealand share trading app Sharesies. Sharesies is a brokerage account, and as such, is a foreign financial account that must be reported to the IRS. FBARs require that you submit the high-water mark for any financial account that you held during the year. The high-water mark is the highest value of the account during the year. Sharesies doesn't provide this information, so you have to calculate it yourself (usually by hovering over the little graph in the portfolio screen and squinting real hard). This program calculates portfolio high-water marks for the given year in an automated and hopefully accurate way.

## Usage

Ensure you have the go toolchain installed, and then clone this repo and run `go run main.go`, and then follow the prompts.

> **Warning** This tool will ask you for your Sharesies login details, and could be doing any number of things. Ensure that you've read the code, and ensured that it's not doing anything nasty. This program could be (but isn't) doing all sorts of nasty things with that information that would be difficult to reverse.

## Contributing

PRs and issues most welcome! I don't expect this to be high enough traffic to need any rules, so it's contribution anarchy over here until that changes.
