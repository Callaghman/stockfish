package internal

import (
	"fmt"
	"io"
	"os/exec"
)

type Stockfish struct {
    Path                    string
    Depth                   int
    Parameters              map[string]interface{}
    NumNodes                int
    TurnPerspective         bool 
    DebugView               bool
    QuitCommandSent         bool
    Info                    string
    DefaultStockfishParams  map[string]interface{}
    StockfishCmd            *exec.Cmd
    Stdin                   io.WriteCloser
    Stdout                  io.ReadCloser
}

func NewStockfish(
    path                string, 
    depth               int, 
    parameters          map[string]interface{},
    numNodes            int, 
    turnPerspective     bool,
    debugView           bool,
) (*Stockfish, error) {
    if path == "" {
        return nil, fmt.Errorf("path to stockfish binary cannot be empty")
    }
    
    if parameters == nil {
        parameters = DefaultStockfishParams
    }

    cmd := exec.Command(path)
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, err
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, err
    }

    if err := cmd.Start(); err != nil {
        return nil, err
    }

    stockfish := &Stockfish{
        Path:               path, 
        Depth:              depth,
        Parameters:         parameters,
        NumNodes:           numNodes,
        TurnPerspective:    turnPerspective,
        DebugView:          debugView,
        QuitCommandSent:    false,
        StockfishCmd:       cmd,
        Stdin:              stdin,
        Stdout:             stdout,
    }

    if err := stockfish.put("uci"); err != nil {
        return nil, err
    }

    stockfish.SetDepth(depth)
    stockfish.SetNumNodes(numNodes)
    stockfish.SetTurnPerspective(turnPerspective)
    stockfish.UpdateEngineParameters(stockfish.Parameters)

    if stockfish.EngineHasWDLOption(){
        stockfish.SetOption("UCI_ShowWDL", true)
    }

    stockfish.PrepareForNewPosition(true)

    return stockfish, nil
}

func (s *Stockfish) SetDebugView(activate bool) {
    s.DebugView = activate
}

func (s *Stockfish) GetEngineParameters() (map[string]interface{}, error) {

    copyParams, err := s.copyEngineParameters() 

    if err != nil {
        return nil, err
    }

    return copyParams, nil
}


func (s *Stockfish) updateEngineParameters(parameters map[string]interface{}) error {
    if len(parameters) == 0 {
        return nil
    }

    newParamValues, err := s.copyEngineParameters()
    
    if err != nil {
        return err
    }

    for k, v := range parameters {
        newParamValues[k] = v
    }

    for key := range newParamValues {
        if len(s.Parameters) > 0 && s.Parameters[key] != nil {
            return fmt.Errorf("Key Error: %s does not exist", key)
        }

        if err := s.validateParamValue(key, newParamValues[key]); err != nil {
            return err
        }
    }
    _, hasSkillLevel := newParamValues["Skill Level"]
    _, hasUCIElo:= newParamValues["UCI_Elo"]
    _, hasUCILimitStrength:= newParamValues["UCI_LimitStrength"]
    
    if(hasSkillLevel != hasUCIElo) && (!hasUCILimitStrength) {
        if hasSkillLevel {
            newParamValues["UCI_LimitStrength"] = false
        } else if hasUCIElo {
            newParamValues["UCI_LimitStrength"] = true
        }
    }

    threadsValue, hasThreads := newParamValues["Threads"]
    var hashValue interface{}
    if hasThreads {
        delete(newParamValues, "Threads")

        if val, hasHash := newParamValues["Hash"]; hasHash {
            hashValue = val
            delete(newParamValues, "Hash")
        } else {
            hashValue = s.Parameters["Hash"]
        }

        newParamValues["Threads"] = threadsValue
        newParamValues["Hash"] = hashValue
    }

    for name, value := range newParamValues {
        if err := s.SetOption(name, value); err != nil {
            return err
        }
    }

    // Set the FEN position after updating UCI options
    if err := s.SetFENPosition(s.GetFENPosition(), false); err != nil {
        return err
    }



    return nil
}

func (s *Stockfish) validateParamValue(key string, value interface{}) error {
    restriction, exists := ParamRestrictions[key]
    if !exists {
        return fmt.Errorf("Parameter Error: Unknown parameter %s", key)
    }

    switch restriction.Type {
    case "int":
        v, ok := value.(int)
        if !ok {
            return fmt.Errorf("Type Error: Invalid type for key %s. Expected int, got %T", key, value)
        }

        if restriction.Min != nil && restriction.Min.(int) > v{
            return fmt.Errorf("Argument Error: Value %v is too small Please select a value in the range (%s, %s)", v, restriction.Min, restriction.Max) 
        }

        if restriction.Min != nil && restriction.Max.(int) < v{
            return fmt.Errorf("Argument Error: Value %v is too large. Please select a value in the range (%s, %s)", v, restriction.Min, restriction.Max)
        }
    case "bool":
        _, ok := value.(bool)
        if !ok {
            return fmt.Errorf("Type Error: Invalid type for key %s: expected bool, got %T", key, value)
        }
    case "string":
        _, ok := value.(string)
        if !ok {
            return fmt.Errorf("Type Error: Invalid type for key %s: expected string, got %T", key, value)
        }
    default: 
        return fmt.Errorf("Type Error: Unsupported type for key %s and value %s. Received type %T", key, value, value) 
    }

    return nil
}

func (s *Stockfish) copyEngineParameters() (map[string]interface{}, error) {

    copyMap := make(map[string]interface{})

    for key, value := range s.Parameters {
        copyMap[key] = value
        switch v := value.(type) {
        case int, bool, string:
            copyMap[key] = v
        default:
            return nil, fmt.Errorf("Error: unsupported type for key %s, value %v: %T\n", key, value, v)
        }
    }
    return copyMap, nil
}

func (s *Stockfish) put(command string) error {
    _, err := fmt.Fprint(s.Stdin, command)
    return err
}
