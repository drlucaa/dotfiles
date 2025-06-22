function jcmb --description "JJ Commit and move/create a jj bookmark"
    # 1. Check for the correct number of arguments
    if test (count $argv) -ne 2
        echo "Usage: jcmb <bookmark_name> <commit_message>"
        echo "Example: jcmb my-feature 'feat: add new button'"
        return 1
    end

    set bookmark_name $argv[1]
    set commit_message $argv[2]

    # 2. First, attempt the commit. If it fails, stop.
    if not jj commit -m "$commit_message"
        echo "jj commit failed. Aborting bookmark operation."
        return 1
    end

    # 3. On successful commit, check if the bookmark exists.
    if jj bookmark list | grep -q "^$bookmark_name"
        # 4a. If it exists, move it to the new commit (@-).
        echo "Moving bookmark '$bookmark_name' to new commit."
        jj bookmark move "$bookmark_name" --to @-
    else
        # 4b. If it doesn't exist, create it at the new commit (@-).
        echo "Creating bookmark '$bookmark_name' at new commit."
        jj bookmark create "$bookmark_name" --revision @-
    end
end
