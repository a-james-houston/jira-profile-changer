# update-jira-profile
updates your jira profile pic from a selected list

you gotta upload the pics first then use the ids

copy env.json.EXAMPLE to env.json and populate with the fields

"avatar_ids" field is a list of int ids of the profile picture ids you want the program to set

"base_url" url to jira instance, e.g. jira.companyname.tld

"usage_order" is the order you want the program to set the pictures in. Supported orders are "SEQUENTIAL" and "RANDOM"

"data_file" is a text file used for storing data between runs. Default data.txt is fine to leave as-is.