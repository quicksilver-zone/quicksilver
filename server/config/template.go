package config

// DefaultConfigTemplate defines the configuration template.
const DefaultConfigTemplate = `
###############################################################################
###                             TLS Configuration                           ###
###############################################################################
[tls]
# Certificate path defines the cert.pem file path for the TLS configuration.
certificate-path = "{{ .TLS.CertificatePath }}"
# Key path defines the key.pem file path for the TLS configuration.
key-path = "{{ .TLS.KeyPath }}"
`
