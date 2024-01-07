package qpm

import "testing"

func TestQPM_job_shell(t *testing.T) {
	t.Parallel()

	type want struct {
		shell string
		ok    bool
	}
	tests := map[string]struct {
		shell map[string]struct{}
		arg   []string
		want  want
	}{
		"argが空の場合、無効": {
			map[string]struct{}{"zsh": {}}, []string{}, want{"", false},
		},
		"shellが空の場合、無効": {
			map[string]struct{}{}, []string{"zsh"}, want{"", false},
		},
		"argがshellに存在しない場合、無効": {
			map[string]struct{}{"zsh": {}}, []string{"bash"}, want{"", false},
		},
		"argがshellに存在する場合、有効": {
			map[string]struct{}{"zsh": {}}, []string{"zsh"}, want{"zsh", true},
		},
		"有効なshellが複数ある場合、優先度の高いものが返る1": {
			map[string]struct{}{"zsh": {}, "bash": {}}, []string{"zsh", "bash"}, want{"zsh", true},
		},
		"有効なshellが複数ある場合、優先度の高いものが返る2": {
			map[string]struct{}{"zsh": {}, "bash": {}}, []string{"bash", "zsh"}, want{"bash", true},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			j := job{
				availableShell: tt.shell,
			}

			// Act
			gotShell, gotOK := j.shell(tt.arg)

			// Assert
			if gotShell != tt.want.shell {
				t.Errorf("job.AvailableShell() got = %v, want %v", gotShell, tt.want.shell)
			}
			if gotOK != tt.want.ok {
				t.Errorf("job.AvailableShell() got1 = %v, want %v", gotOK, tt.want.ok)
			}
		})
	}
}
