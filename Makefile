app:
	gcloud beta run deploy elos-app --project spin-309720 --service-account elos-app-identity --platform managed --region us-west1 --source .
mod:
	GOPRIVATE=github.com/nlandolfi,github.com/spinsrv go mod tidy
