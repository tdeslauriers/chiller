# Chiller

Support service to back up the [deslauriers.world](https://deslauriers.world) site data.

## Context:

Because the site is hosted on broken computers, and more importantly, because I am constantly fiddling with it, crashes and failures happen.  

*Learned some hard lessons I already knew about why external backups are important.*

## Function: Runs nigthly back ups

1. Logs into site. 
1. Uses read access to micro-services to pull db table data.
1. Loads and encrypts data into back up database outside K8s cluster.
    * field level encryption.

## Day Two

1.  Build function to populate data back to active site in the event of cluster failure.  