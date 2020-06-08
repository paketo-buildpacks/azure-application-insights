printf "Configuring Azure Application Insights properties\n"

# shellcheck disable=SC2046
eval $(azure-application-insights-properties) || exit $?
