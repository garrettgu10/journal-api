# ggu's Journal API

The API will log journal entries to the correct file in a Git repository. 

The file structure of the journal repository will be as follows:
```
/
    2022/
        08/
            05
            06
            ...
        07/
            ...
```

Only two features are supported:
- Add a new entry -- this accepts the current date & an entry. Creates a new line in the corresponding file and appends the entry.
- Commit changes -- this adds all pending changes to the git repository and pushes the changes. 

This API does not allow you to read entries, only create new ones.

## Deploying App
Make a `.env` file based on the provided `.env_example` file, and deploy the app using `./run.sh`. I also put mine behind an nginx reverse proxy so I can use HTTPS easily. 

## Building `/static`

The /static directory was added for deployment convenience & it's just a copy of a build from [https://github.com/garrettgu10/journal-frontend](https://github.com/garrettgu10/journal-frontend)

You can update it by doing the following: (assuming journal-frontend and journal-api are under the same directory)
```
cd journal-frontend
npm run-script build
mv build ../journal-api/static
```
