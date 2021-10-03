# home-assistant-to-big-query

## Why ?

Thanks to this project, you will link your the data from your [HomeAssistant](https://www.home-assistant.io/) to [BigQuery](https://cloud.google.com/bigquery)

## Prerequisite

- a [GCP account](https://console.cloud.google.com)
- Home assistant and its [PubSub plugin](https://www.home-assistant.io/integrations/google_pubsub/)

## Quick start

Deploy this project as a [Google Cloud Function](https://cloud.google.com/functions) in your GCP Project. This function is triggered by the event pushed by home assistant.
It needs 3 environment variables : 
- **GCP_PROJECTID**
- **BIGQUERY_DATASET** 
- **BIGQUERY_TABLE**

## Tutorial

A complete tutorial is available at https://blog.bodul.fr