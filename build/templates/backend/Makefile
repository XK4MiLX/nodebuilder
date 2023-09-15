{{define "main" -}}
ARCHIVE := $(shell basename {{.Backend.BinaryURL}})

all:
	if [ -d "backend" ]; then rm -rf backend; fi
	wget --tries 5 {{.Backend.BinaryURL}} -q --show-progress
{{- if eq .Backend.VerificationType "gpg"}}
	wget {{.Backend.VerificationSource}} -O checksum
	gpg --batch --import ../../build/gpg-keys/*
	gpg --verify checksum ${ARCHIVE}
{{- else if eq .Backend.VerificationType "gpg-sha256"}}
	wget {{.Backend.VerificationSource}} -O checksum
	gpg --batch --import ../../build/gpg-keys/*
	gpg --verify checksum
	sha256sum -c --ignore-missing checksum
{{- else if eq .Backend.VerificationType "sha256"}}
	[ "$$(sha256sum ${ARCHIVE} | cut -d ' ' -f 1)" = "{{.Backend.VerificationSource}}" ]
{{- end}}
	mkdir backend
	{{.Backend.ExtractCommand}} ${ARCHIVE}
        chmod 755 debian/bootstrap
	chmod 755 debian/corruption
	chmod 755 debian/health
	chmod 755 debian/setup
	chmod 755 debian/postinst
	chmod 755 debian/run
{{- if .Backend.ExcludeFiles}}
	# generated from exclude_files
{{- range $index, $name := .Backend.ExcludeFiles}}
	rm backend/{{$name}}
{{- end}}
{{- end}}
	bash -c "debian/setup"
clean:
	rm -rf backend
	rm -f ${ARCHIVE}
	rm -f checksum
{{end}}
