## TicTacToe2

### Build

``` sh
make build
```

executable `tt` appears

### Usage 

``` sh
tt [train|playX|playO]
```

`tt train` will train the model by playing it against a minimax player and then by another reinforcement learning player. This will generate the file `learner_player.json`

`tt playX` will play the model against a human player. The human player will play first as X. It is expected that the model file `learner_player.json` has been adequately trained.

`tt playO` same as `tt playX` except the human play will play second as O.

