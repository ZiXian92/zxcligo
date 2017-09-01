package zxcligo

import (
	"reflect"
	"testing"
)

func Test_newContext(t *testing.T) {
	type args struct {
		cmdStrings []string
	}
	tests := []struct {
		name    string
		args    args
		wantCtx *Context
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCtx, err := newContext(tt.args.cmdStrings)
			if (err != nil) != tt.wantErr {
				t.Errorf("newContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCtx, tt.wantCtx) {
				t.Errorf("newContext() = %v, want %v", gotCtx, tt.wantCtx)
			}
		})
	}
}

func Test_isLongNameOption(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Option name only",
			args: args{s: "--name"},
			want: true,
		},
		{
			name: "Option name with value",
			args: args{s: "--name=andy"},
			want: true,
		},
		{
			name: "Short name option",
			args: args{s: "---name=andy"},
			want: false,
		},
		{
			name: "No -- prefix",
			args: args{s: "name"},
			want: false,
		},
		{
			name: "Triple dash",
			args: args{s: "---"},
			want: false,
		},
		{
			name: "Correct prefix no opt name",
			args: args{s: "--"},
			want: false,
		},
		{
			name: "Single dash",
			args: args{s: "-"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLongNameOption(tt.args.s); got != tt.want {
				t.Errorf("isLongNameOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isShortNameOption(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Short name only",
			args: args{s: "-n"},
			want: true,
		},
		{
			name: "Short name with value",
			args: args{s: "-n=andy"},
			want: true,
		},
		{
			name: "Long name only",
			args: args{s: "--name"},
			want: false,
		},
		{
			name: "Long name with value",
			args: args{s: "--name=andy"},
			want: false,
		},
		{
			name: "No prefix",
			args: args{s: "n"},
			want: false,
		},
		{
			name: "Double dash",
			args: args{s: "--"},
			want: false,
		},
		{
			name: "Single dash",
			args: args{s: "-"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isShortNameOption(tt.args.s); got != tt.want {
				t.Errorf("isShortNameOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOption(t *testing.T) {
	type args struct {
		optString string
	}
	tests := []struct {
		name         string
		args         args
		wantOptName  string
		wantOptValue string
		wantErr      bool
	}{
		{
			name:         "Long name option",
			args:         args{optString: "--name"},
			wantOptName:  "name",
			wantOptValue: "",
			wantErr:      false,
		},
		{
			name:         "Long name option with value",
			args:         args{optString: "--name=andy"},
			wantOptName:  "name",
			wantOptValue: "andy",
			wantErr:      false,
		},
		{
			name:         "Short name option",
			args:         args{optString: "-n"},
			wantOptName:  "n",
			wantOptValue: "",
			wantErr:      false,
		},
		{
			name:         "Short name option with value",
			args:         args{optString: "-n=andy"},
			wantOptName:  "n",
			wantOptValue: "andy",
			wantErr:      false,
		},
		{
			name:         "Triple dash prefix",
			args:         args{optString: "---name"},
			wantOptName:  "",
			wantOptValue: "",
			wantErr:      true,
		},
		{
			name:         "Normal argument",
			args:         args{optString: "name"},
			wantOptName:  "",
			wantOptValue: "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOptName, gotOptValue, err := parseOption(tt.args.optString)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOptName != tt.wantOptName {
				t.Errorf("parseOption() gotOptName = %v, want %v", gotOptName, tt.wantOptName)
			}
			if gotOptValue != tt.wantOptValue {
				t.Errorf("parseOption() gotOptValue = %v, want %v", gotOptValue, tt.wantOptValue)
			}
		})
	}
}

func Test_processCommandOptions(t *testing.T) {
	type args struct {
		ctx     *Context
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Conflicting options",
			args: args{
				ctx: &Context{Options: map[string]interface{}{
					"name": "andy",
					"n":    "andy",
				}},
				options: []Option{
					{
						LongName:  "name",
						ShortName: "n",
						Type:      OptTypeString,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Incorrect value type",
			args: args{
				ctx: &Context{Options: map[string]interface{}{
					"numtries": "andy",
				}},
				options: []Option{
					{
						LongName: "numTries",
						Type:     OptTypeInt,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Unsupported options",
			args: args{
				ctx: &Context{Options: map[string]interface{}{
					"message": "You fool!",
					"name":    "andy",
				}},
				options: []Option{
					{
						LongName: "andy",
						Type:     OptTypeString,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Long and short name same",
			args: args{
				ctx: &Context{Options: map[string]interface{}{
					"age": "27",
				}},
				options: []Option{
					{
						LongName:     "age",
						ShortName:    "age",
						Type:         OptTypeUint,
						DefaultValue: 23,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No default value",
			args: args{
				ctx: &Context{Options: map[string]interface{}{}},
				options: []Option{
					{
						LongName: "name",
						Type:     OptTypeString,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Ok case",
			args: args{
				ctx: &Context{Options: map[string]interface{}{
					"age": "23",
				}},
				options: []Option{
					{
						LongName:     "name",
						Type:         OptTypeString,
						DefaultValue: "andy",
					},
					{
						ShortName: "age",
						Type:      OptTypeUint,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := processCommandOptions(tt.args.ctx, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("processCommandOptions() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				for _, opt := range tt.args.options {
					var name string
					if len(opt.LongName) > 0 {
						name = opt.LongName
					} else {
						name = opt.ShortName
					}
					_, ok := tt.args.ctx.Options[name]
					if !ok {
						t.Errorf("Option %s's value not present!", name)
					}
				}
			}
		})
	}
}
