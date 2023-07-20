package stable_diffusion

type Extension interface {
	Key() string
	Value() map[string]any
}

type Extensions struct {
	Controlnet *Controlnet `json:"controlnet,omitempty"`
	// more...
}

func (exts *Extensions) HasControlnetImg() bool {
	if exts == nil {
		return false
	}

	if exts.Controlnet == nil {
		return false
	}

	if exts.Controlnet.InputImage == "" {
		return false
	}

	return true
}

type Controlnet struct {
	InputImage string  `json:"input_image,omitempty" mapstructure:"input_image"` // default=""
	Mask       string  `json:"mask,omitempty" mapstructure:"mask"`               // default=""
	Module     string  `json:"module,omitempty" mapstructure:"module"`           // default="none"
	Model      string  `json:"model,omitempty" mapstructure:"model"`             // default="none"
	Guessmode  bool    `json:"guessmode,omitempty" mapstructure:"guessmode"`     // default=true
	Lowvram    bool    `json:"lowvram,omitempty" mapstructure:"lowvram"`         // default: false
	Weight     float64 `json:"weight,omitempty" mapstructure:"weight"`           // default=1.0
	// enum: 0=Just Resize, 1=Scale to Fit (Inner Fit), 2=Envelope (Outer Fit)
	// default=1, "Scale to Fit (Inner Fit)"
	ResizeMode    uint    `json:"resize_mode,omitempty" mapstructure:"resize_mode"`
	ProcessorRes  int     `json:"enabled,omitempty" mapstructure:"resize_mode"`           // default=64
	GuidanceStart float64 `json:"guidance_start,omitempty" mapstructure:"guidance_start"` // default=0.0
	GuidanceEnd   float64 `json:"guidance_end,omitempty" mapstructure:"guidance_end"`     // default=1.0
	// ThresholdA    float64 `json:"threshold_a,omitempty"`    // default=64
	// ThresholdB    float64 `json:"threshold_b,omitempty"`    // default=64
}

/*
module	model
canny	control_sd15_canny [fe5fe48e]
depth	control_sd15_depth [fe5fe48e]
mlsd	control_sd15_mlsd [fe5fe48e]
*/

func (item *Controlnet) Key() string {
	return "controlnet"
}

func (item *Controlnet) Value() map[string]any {
	return map[string]any{"args": []any{item}}
}

func (exts *Extensions) List() (results []Extension) {
	results = make([]Extension, 0)

	if exts.Controlnet != nil && exts.Controlnet.InputImage != "" {
		results = append(results, exts.Controlnet)
	}

	return results
}
