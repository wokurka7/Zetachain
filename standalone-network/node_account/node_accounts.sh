
zetacored tx observer add-observer zeta1cxj07f3ju484ry2cnnhxl5tryyex7gev0yzxtj zetapub1addwnpepqwwpjwwnes7cywfkr0afme7ymk8rf5jzhn8pfr6qqvfm9v342486qsrh4f5 true --from zeta --gas=auto --gas-prices=10000000000azeta --gas-adjustment=1.5 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx observer add-observer zeta1v66xndg92tkt9ay70yyj3udaq0ej9wl765r7lf zetapub1addwnpepqwwpjwwnes7cywfkr0afme7ymk8rf5jzhn8pfr6qqvfm9v342486qsrh4f5 false --from zeta --gas=auto --gas-prices=10000000000azeta --gas-adjustment=1.5 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block


zetacored q observer show-node-account zeta1cxj07f3ju484ry2cnnhxl5tryyex7gev0yzxtj
zetacored q observer show-node-account zeta1v66xndg92tkt9ay70yyj3udaq0ej9wl765r7lf
zetacored q observer list-observer-set

zetacored tx observer add-observer zeta1v66xndg92tkt9ay70yyj3udaq0ej9wl765r7lf zetapub1addwnpepqwwpjwwnes7cywfkr0afme7ymk8rf5jzhn8pfr6qqvfm9v342486qsrh4f5 true --from zeta --gas=auto --gas-prices=10000000000azeta --gas-adjustment=1.5 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx observer update-observer zeta1v66xndg92tkt9ay70yyj3udaq0ej9wl765r7lf zeta1cxj07f3ju484ry2cnnhxl5tryyex7gev0yzxtj 2 --from zeta --gas=auto --gas-prices=10000000000azeta --gas-adjustment=1.5 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block

