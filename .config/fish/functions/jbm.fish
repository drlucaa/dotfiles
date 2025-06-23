function jbm --description "JJ move/create bookmark to/at current change"
    # 1. Check for the correct number of arguments
    if test (count $argv) -ne 1
        echo "Usage: jbm <bookmark_name>"
        echo "Example: jbm my-feature"
        return 1
    end

    set bookmark_name $argv[1]

    # 2. Check if the bookmark exists.
    if jj bookmark list | grep -q "^$bookmark_name"
        # 3a. If it exists, move it to the current change (@).
        echo "Moving bookmark '$bookmark_name' to new change."
        jj bookmark move "$bookmark_name" --to @
    else
        # 3b. If it doesn't exist, create it at the current change (@).
        echo "Creating bookmark '$bookmark_name' at new change."
        jj bookmark create "$bookmark_name" --revision @
    end
end
