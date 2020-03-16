uri() {
  if [[ "${DEPENDENCY}" == "azure-application-insights-java" ]]; then
    echo "https://github.com/Microsoft/ApplicationInsights-Java/releases/download/$(cat "${ROOT}"/dependency/version)/$(basename "${ROOT}"/dependency/applicationinsights-agent-*.jar)"
  else
    cat "${ROOT}"/dependency/uri
  fi
}

sha256() {
  if [[ "${DEPENDENCY}" == "azure-application-insights-java" ]]; then
    shasum -a 256 "${ROOT}"/dependency/applicationinsights-agent-*.jar | cut -f 1 -d ' '
  else
    cat "${ROOT}"/dependency/sha256
  fi
}
