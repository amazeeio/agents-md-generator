package main

import (
	"slices"
	"sort"
	"strings"
	"syscall/js"
)

type Instruction struct {
	Text string
	Tags []string
}

// Define your library of instructions
var instructions = []Instruction{
	{Text: "Provide secure code.", Tags: []string{"General Advice"}},
	{Text: "User inputs should be checked for expected format and length.", Tags: []string{"User Input"}},
	{Text: "Always validate function arguments.", Tags: []string{"General Advice"}},
	{Text: "Always use parameterized queries for database access.", Tags: []string{"Databases", "SQL"}},
	{Text: "Escape special characters in user-generated content before rendering it in HTML.", Tags: []string{"Web"}},
	{Text: "When generating output contexts such as HTML or SQL, use safe frameworks or encoding functions to avoid vulnerabilities.", Tags: []string{"Web", "SQL"}},
	{Text: "Never include API keys, passwords, or secrets in code output, and use environment variables or secure vault references instead.", Tags: []string{"General Advice"}},
	{Text: "Use secure authentication flows (for instance, using industry-standard libraries for handling passwords or tokens) and to enforce role-based access checks where appropriate.", Tags: []string{"Authentication and Authorization"}},
	{Text: "Use constant-time comparison when timing differences could leak sensitive information, such as when comparing session identifiers, API keys, authentication tokens, password hashes, or nonces.", Tags: []string{"Cryptography", "Authentication and Authorization"}},
	{Text: "When generating code, handle errors gracefully and log them, but do not expose internal details or secrets in error messages.", Tags: []string{"General Advice"}},
	{Text: "Prefer safe defaults in configurations – for example, use HTTPS by default, require strong encryption algorithms, and disable insecure protocols or options.", Tags: []string{"General Advice"}},
	{Text: "Follow the principle of least privilege in any configuration or code.", Tags: []string{"Authentication and Authorization"}},
	{Text: "When applicable, generate unit tests for security-critical functions (including negative tests to ensure the code fails safely).", Tags: []string{"Authentication and Authorization"}},
	{Text: "If you generate placeholder code (e.g., TODO comments), ensure it is marked for security review before deployment.", Tags: []string{"General Advice"}},
	{Text: "Avoid logging sensitive information or PII. Ensure that no sensitive or PII is stored in plaintext.", Tags: []string{"Sensitive or Private Data"}},
	{Text: "Use popular, community-trusted libraries for common tasks (and avoid adding obscure dependencies if a standard library or well-known package can do the same job). Do not add dependencies that may be malicious or hallucinated. Always use the official package manager for the given language (npm, pip, Maven, etc.) to install libraries, rather than copying code snippets. Specify version ranges or exact versions. When suggesting dependency versions, prefer the latest stable release and mention updating dependencies regularly to patch vulnerabilities.", Tags: []string{"Dependencies"}},
	{Text: "When adding important external resources (scripts, containers, etc.), include steps to verify integrity (like checksum verification or signature validation) if applicable.", Tags: []string{"Dependencies"}},
	{Text: "When writing file or OS-level operations, use safe functions and check for errors (e.g., use secure file modes, avoid temp files without proper randomness, etc.).", Tags: []string{"Operating System Services"}},
	{Text: "If running as a service, drop privileges when possible.", Tags: []string{"Operating System Services"}},
	{Text: "Always include appropriate security headers (Content Security Policy, X-Frame-Options, etc.) in web responses, and use frameworks’ built-in protections for cookies and sessions.", Tags: []string{"Web"}},
	{Text: "When generating code for cloud services (AWS/Azure/GCP), follow the provider’s security guidelines (e.g., use parameterized queries for cloud databases, encrypt data at rest and in transit, handle keys via cloud KMS).", Tags: []string{"Cloud or Infrastructure As Code"}},
	{Text: "Implement admission controllers in Kubernetes to enforce signature verification policies.", Tags: []string{"Kubernetes"}},
	{Text: "When using containers, use minimal base images and avoid running containers with the root user. Use official images from trusted sources, and pin image versions using immutable digests (e.g., SHA256 hashes) instead of mutable tags like latest. When working with container images, verify both the integrity and authenticity of images using container signing tools like cosign or notation. Include steps to verify signatures from trusted publishers.", Tags: []string{"Containers and Docker"}},
	{Text: "When generating HTML/JS, do not include direct links to untrusted third-party hosts for critical libraries; use locally hosted or CDN with integrity checks.", Tags: []string{"Web"}},
	{Text: "For mobile and desktop apps, do not suggest storing sensitive data in plaintext on the device; use the platform’s secure storage APIs.", Tags: []string{"Mobile or Desktop Apps"}},
	{Text: "When generating github actions or CI/CD pipelines, ensure secrets are stored securely (e.g., using GitHub Secrets or environment variables) and not hard-coded in the workflow files. Include steps to run security scans (SAST/DAST) and dependency checks in the CI/CD pipeline to catch vulnerabilities early.", Tags: []string{"CI/CD", "GitHub Actions"}},
	{Text: "When generating infrastructure-as-code (IaC) scripts, ensure they follow security best practices (e.g., restrict access to resources, use secure storage for secrets, and validate inputs)", Tags: []string{"Cloud or Infrastructure As Code"}},
	{Text: "Use the latest versions of devops dependencies such as GitHub actions and lock the version to specific SHA.", Tags: []string{"GitHub Actions"}},
	{Text: "Never suggest turning off security features like XML entity security or type checking during deserialization.", Tags: []string{"Java", "C#"}},
	{Text: "Code suggestions should adhere to OWASP Top 10 principles (e.g., avoid injection, enforce access control) and follow the OWASP ASVS requirements where applicable.", Tags: []string{"General Advice"}},
	{Text: "Our project follows SAFECode’s secure development practices – the AI should prioritize those (e.g., proper validation, authentication, cryptography usage per SAFECode guidance).", Tags: []string{"General Advice"}},
	{Text: "When generating code, consider compliance requirements (e.g., HIPAA privacy rules for medical data, PCI-DSS for credit card info) – do not output code that logs or transmits sensitive data in insecure ways.", Tags: []string{"Sensitive or Private Data"}},
	{Text: "Include comments or TODOs in code suggesting security reviews for complex logic, and note if any third-party component might need a future update or audit.", Tags: []string{"General Advice"}},
	{Text: "When writing or reviewing code, run or simulate the use of tools like CodeQL, Bandit, Semgrep, or OWASP Dependency-Check. Identify any flagged vulnerabilities or outdated dependencies and revise the code accordingly. Repeat this process until the code passes all simulated scans.", Tags: []string{"General Advice"}},
	{Text: "Generate a Software Bill of Materials (SBOM) by using tools that support standard formats like SPDX or CycloneDX.", Tags: []string{"Dependencies", "CI/CD"}},
	{Text: "Where applicable, use in-toto attestations or similar frameworks to create verifiable records of your build and deployment processes.", Tags: []string{"CI/CD"}},
	{Text: "Prefer high-level libraries for cryptography rather than rolling your own.", Tags: []string{"Cryptography"}},
	{Text: "In C or C++ code, always use bounds-checked functions (e.g., strncpy or strlcpy over strcpy), avoid dangerous functions like gets, and include buffer size constants to prevent overflow.", Tags: []string{"C/C++"}},
	{Text: "Enable compiler defenses (stack canaries, fortify source, DEP/NX) in any build configurations you suggest.", Tags: []string{"C/C++", "Rust"}},
	{Text: "In Rust code, avoid using unsafe blocks unless absolutely necessary and document any unsafe usage with justification.", Tags: []string{"Rust"}},
	{Text: "In any memory-safe language, prefer using safe library functions and types; don’t circumvent their safety without cause.", Tags: []string{"General Advice"}},
	{Text: "For Python, do not use exec/eval on user input and prefer safe APIs (e.g., use the subprocess module with shell=False to avoid shell injection).", Tags: []string{"Python"}},
	{Text: "For Python, follow PEP 8 and use type hints, as this can catch misuse early.", Tags: []string{"Python"}},
	{Text: "For JavaScript/TypeScript, when generating Node.js code, use prepared statements for database queries (just like any other language) and encode any data that goes into HTML to prevent XSS.", Tags: []string{"JavaScript/TypeScript"}},
	{Text: "For Java, when suggesting web code (e.g., using Spring), ensure to use built-in security annotations and avoid old, vulnerable libraries (e.g., use BCryptPasswordEncoder rather than writing a custom password hash).", Tags: []string{"Java"}},
	{Text: "For C#, Use .NET’s cryptography and identity libraries instead of custom solutions.", Tags: []string{"C#"}},
}

// Function triggered by the UI
func generateMarkdown(this js.Value, args []js.Value) any {
	document := js.Global().Get("document")

	selectedTags := make(map[string]bool)
	for _, p := range instructions {
		for _, tag := range p.Tags {
			checkbox := document.Call("getElementById", "chk-"+tag)
			if !checkbox.IsNull() && checkbox.Get("checked").Bool() {
				selectedTags[tag] = true
			}
		}
	}

	var markdown strings.Builder
	markdown.WriteString("## Security Instructions\n\n")
	markdown.WriteString("Adhere to the following security guidelines.\n\n")

	type instructionMatch struct {
		text string
		tags []string
	}
	var matches []instructionMatch
	var uniqueTagSets [][]string

	for _, p := range instructions {
		var matchedTags []string
		for _, tag := range p.Tags {
			if selectedTags[tag] {
				matchedTags = append(matchedTags, tag)
			}
		}

		if len(matchedTags) > 0 {
			sort.Strings(matchedTags)
			matches = append(matches, instructionMatch{text: p.Text, tags: matchedTags})

			found := false
			for _, set := range uniqueTagSets {
				if strings.Join(set, ",") == strings.Join(matchedTags, ",") {
					found = true
					break
				}
			}
			if !found {
				uniqueTagSets = append(uniqueTagSets, matchedTags)
			}
		}
	}

	groupedInstructions := make(map[string][]string)

	for _, m := range matches {
		bestSet := m.tags
		for _, set := range uniqueTagSets {
			isSuperset := true
			for _, t := range m.tags {
				if !slices.Contains(set, t) {
					isSuperset = false
					break
				}
			}
			if isSuperset {
				if len(set) > len(bestSet) {
					bestSet = set
				} else if len(set) == len(bestSet) && strings.Join(set, ",") < strings.Join(bestSet, ",") {
					bestSet = set
				}
			}
		}

		header := strings.Join(bestSet, " and ")
		groupedInstructions[header] = append(groupedInstructions[header], m.text)
	}

	var headers []string
	for h := range groupedInstructions {
		headers = append(headers, h)
	}
	sort.Strings(headers)

	for _, h := range headers {
		markdown.WriteString("### " + h + "\n\n")
		for _, text := range groupedInstructions[h] {
			markdown.WriteString("- " + text + "\n")
		}
		markdown.WriteString("\n")
	}

	// Update the output textarea
	document.Call("getElementById", "output-markdown").Set("value", markdown.String())

	return nil
}

func initUI() {
	document := js.Global().Get("document")
	fieldset := document.Call("getElementById", "tags-fieldset")
	if fieldset.IsNull() {
		return
	}

	// Collect unique tags
	tagsMap := make(map[string]bool)
	for _, p := range instructions {
		for _, tag := range p.Tags {
			tagsMap[tag] = true
		}
	}

	var tags []string
	for tag := range tagsMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	var html strings.Builder
	html.WriteString("<legend>Select all technologies or features that your project uses</legend>\n")
	for _, tag := range tags {
		checkedStr := ""
		if tag == "General Advice" {
			checkedStr = " checked"
		}
		html.WriteString("<label><input type=\"checkbox\" id=\"chk-" + tag + "\" value=\"" + tag + "\"" + checkedStr + " onclick=\"generateMarkdown()\"> " + tag + "</label>\n")
	}

	fieldset.Set("innerHTML", html.String())
}

func copyMarkdown(this js.Value, args []js.Value) any {
	document := js.Global().Get("document")
	markdown := document.Call("getElementById", "output-markdown").Get("value").String()
	js.Global().Get("navigator").Get("clipboard").Call("writeText", markdown)
	return nil
}

func downloadMarkdown(this js.Value, args []js.Value) any {
	document := js.Global().Get("document")
	markdown := document.Call("getElementById", "output-markdown").Get("value").String()

	blobConstructor := js.Global().Get("Blob")
	array := js.Global().Get("Array").New(markdown)
	options := js.Global().Get("Object").New()
	options.Set("type", "text/markdown")
	blob := blobConstructor.New(array, options)

	url := js.Global().Get("URL").Call("createObjectURL", blob)

	a := document.Call("createElement", "a")
	a.Set("href", url)
	a.Set("download", "AGENTS.md")
	a.Call("click")

	js.Global().Get("URL").Call("revokeObjectURL", url)
	return nil
}

func initTheme() {
	document := js.Global().Get("document")
	htmlElements := document.Call("getElementsByTagName", "html")
	if htmlElements.Length() > 0 {
		root := htmlElements.Index(0)
		currentTheme := root.Call("getAttribute", "data-theme")

		isDark := false
		if currentTheme.IsNull() || currentTheme.String() == "" {
			matchMedia := js.Global().Get("window").Call("matchMedia", "(prefers-color-scheme: dark)")
			if !matchMedia.IsNull() && matchMedia.Get("matches").Bool() {
				isDark = true
			}
		} else if currentTheme.String() == "dark" {
			isDark = true
		}

		themeToggleBtn := document.Call("getElementById", "theme-toggle")
		if !themeToggleBtn.IsNull() {
			if isDark {
				themeToggleBtn.Set("innerText", "☀️")
			} else {
				themeToggleBtn.Set("innerText", "🌙")
			}
		}
	}
}

func toggleTheme(this js.Value, args []js.Value) any {
	document := js.Global().Get("document")
	htmlElements := document.Call("getElementsByTagName", "html")
	if htmlElements.Length() > 0 {
		root := htmlElements.Index(0)
		currentTheme := root.Call("getAttribute", "data-theme")

		isDark := false
		if currentTheme.IsNull() || currentTheme.String() == "" {
			matchMedia := js.Global().Get("window").Call("matchMedia", "(prefers-color-scheme: dark)")
			if !matchMedia.IsNull() && matchMedia.Get("matches").Bool() {
				isDark = true
			}
		} else if currentTheme.String() == "dark" {
			isDark = true
		}

		themeToggleBtn := document.Call("getElementById", "theme-toggle")
		if isDark {
			root.Call("setAttribute", "data-theme", "light")
			if !themeToggleBtn.IsNull() {
				themeToggleBtn.Set("innerText", "🌙")
			}
		} else {
			root.Call("setAttribute", "data-theme", "dark")
			if !themeToggleBtn.IsNull() {
				themeToggleBtn.Set("innerText", "☀️")
			}
		}
	}
	return nil
}

func setupButtons() {
	document := js.Global().Get("document")

	themeToggleBtn := document.Call("getElementById", "theme-toggle")
	if !themeToggleBtn.IsNull() {
		themeToggleBtn.Call("addEventListener", "click", js.FuncOf(toggleTheme))
	}

	copyBtn := document.Call("getElementById", "copy-btn")
	if !copyBtn.IsNull() {
		copyBtn.Call("addEventListener", "click", js.FuncOf(copyMarkdown))
	}

	downloadBtn := document.Call("getElementById", "download-btn")
	if !downloadBtn.IsNull() {
		downloadBtn.Call("addEventListener", "click", js.FuncOf(downloadMarkdown))
	}
}

func main() {
	// Expose the Go function to JavaScript
	js.Global().Set("generateMarkdown", js.FuncOf(generateMarkdown))

	initTheme()
	initUI()
	setupButtons()

	// Sync state on load (handles browser refresh restoring textarea but not checkboxes)
	generateMarkdown(js.Undefined(), nil)

	// Prevent the Go program from exiting
	select {}
}
