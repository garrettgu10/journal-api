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