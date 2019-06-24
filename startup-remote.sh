echo communicationg with remote $1

curl -X POST $1/q/enqueue -d '{"Name": "Kyle Felter", "Course": "ECE 2003", "Project":"Lab 5", "Email":"felterkyle94@gmail.com"}'