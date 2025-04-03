Learning GO and Ebitengine at the same time.

Game will be:

You are a skeleton. You need to pick flowers to decorate the graves in your cemetery. Heroes sometimes don't understand and try to attack you. If you attack back, you will save yourself, but have the added problem of having MORE graves to add flowers to.

Each level ends when all the flowers are placed. However, in the next level the number of people follow the Fibonacci sequence and grow exponentially, making your caretaking job more difficult.

Run:
`go run main.go`

Run with wasm:
- Create Initial Build with WASM: `go run github.com/hajimehoshi/wasmserve@latest ./path/to/yourgame`
- Then, or subsequent builds: `GOOS=js GOARCH=wasm go build -o skeleton.wasm`      
- Run Server: `go run github.com/hajimehoshi/wasmserve@latest ./main.go`

Controls:
- Move: Arrow keys / touch screen high/low/right/left zone
- Attack: Space / tap
- Toggle Full Screen: F11 or F