function jcbm --description "JJ commit and move/create a bookmark"
    # 1. Check for the correct number of arguments
    if test (count $argv) -ne 2
        echo "Usage: jcbm <bookmark_name> <change_message>"
        echo "Example: jcbm my-feature 'feat: add new button'"
        return 1
    end

    set bookmark_name $argv[1]
    set change_message $argv[2]

    # 2. First, attempt the change. If it fails, stop.
    if not jj commit -m "$change_message"
        echo "jj commit failed. Aborting bookmark operation."
        return 1
    end

    # 3. On successful change, check if the bookmark exists.
    if jj bookmark list | grep -q "^$bookmark_name"
        # 4a. If it exists, move it to the new change (@-).
        echo "Moving bookmark '$bookmark_name' to new change."
        jj bookmark move "$bookmark_name" --to @-
    else
        # 4b. If it doesn't exist, create it at the new change (@-).
        echo "Creating bookmark '$bookmark_name' at new change."
        jj bookmark create "$bookmark_name" --revision @-
    end
end
