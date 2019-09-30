var loadingOverlay; 
$(document).ready(function(){
    loadingOverlay = document.querySelector('.loading');
    // Display the flash message
    window.Flash.create('.flash-message');

    $(".sidebar").slimScroll({
        height: 'auto',
        color: '#CCCCCC'
    });

    // Select a company as active
    // Init a timeout variable to be used below
    var selectcompanytimeout = null;
    var $selectCompanyForm = document.getElementById("select-company-form");    
    $(document).on('click', '.btn-dropdown-select-company', function(){ 
        // Clear the timeout if it has already been set.
        // This will prevent the previous task from executing
        // if it has been less than <MILLISECONDS>
        clearTimeout(selectcompanytimeout);

        var compId = $(this).data("id");
        let csrfToken = document.getElementsByName("gorilla.csrf.Token")[0].value;
        // Make a new timeout set to go off in 800ms
        selectcompanytimeout = setTimeout(function () {
            var selectCompURL = "/dashboard/company/" + compId + "/visit";
            $selectCompanyForm.action = selectCompURL;
            $selectCompanyForm.submit();
            $('#modal-select-company').modal("hide");
            toggleLoading();
            
        }, 500);
    });

    // Search user
    // Get the input box
    var $navbarSearchInput = document.getElementById('input-navbar-search');
    if (typeof($navbarSearchInput) != 'undefined' && $navbarSearchInput != null) {
        // Init a timeout variable to be used below
        var searchTimeout = null;
        var loadingText = '<small class="dropdown-item text-muted"><i class="fa fa-spinner mr-1"></i>Loading...</small>';

        // Listen for keystroke events
        $navbarSearchInput.onkeyup = function (e) {

            // Clear the timeout if it has already been set.
            // This will prevent the previous task from executing
            // if it has been less than <MILLISECONDS>
            clearTimeout(searchTimeout);

            // Make a new timeout set to go off in 1000ms
            searchTimeout = setTimeout(function () {
                var value = $navbarSearchInput.value;
                if(value.trim().length >= 3) {                  
                    var $searchBar = $("#navbar-search-bar");
                    var $searchDropdown = $("#navbar-search-dropdown");

                    // Remove all content in the dropdown first
                    $searchDropdown.empty();

                    // Set the text to be loading
                    $searchDropdown.html(loadingText);

                    // Open the drop down
                    $searchBar.addClass("open");

                    var url = "/dashboard/company/users/search?query=" + value.trim();
                    axios.get(url)
                    .then(function (response) {
                        // handle success 
                        
                        var result = response["data"]["data"];
                        // Remove all content in the dropdown first
                        $searchDropdown.empty();
                        if(response["data"]["success"]) {
                            // Populate the dropdown
                            if(result.length > 0){
                                var limit = 5;
                                var html = '';
                                for(var i = 0; i < limit && i < result.length; i++) {
                                    // Build the search item
                                    html += '<a href="/dashboard/user/' + result[i]["ID"] + '" class="dropdown-item">'
                                            + '<div class="media">'
                                            + '<img src="' + result[i]["profilePicture"] + '" alt="' + result[i]["name"] + '" class="img-size-50 mr-3 img-circle">'
                                            + '<div class="media-body">'
                                            + '<p class="text-sm">' + result[i]["name"] + '</p>'
                                            + '<p class="text-sm text-muted">' + result[i]["email"] + '</p>'
                                            + '</div>'
                                            + '</div>'
                                            + '</a>'
                                            + '<div class="dropdown-divider"></div>';    
                                }
                                html += '<a href="/dashboard/company/users/search/all?query=' + value.trim() + '" class="dropdown-item dropdown-footer">Show all</a>';
                                $searchDropdown.html(html);
                            } else {
                                $searchDropdown.html('<small class="dropdown-item text-muted">No result found.</small>');
                            }
                        } else {
                            $searchDropdown.html('<small class="dropdown-item text-muted">Please ensure that you have selected a company.</small>');
                        }
                        
                    })
                    .catch(function (error) {
                        // handle error
                        console.log(error);
                    });
                }
            }, 1000);
        };

        // Hide the navbar search dropdown if click outside
        $(document).mouseup(function (e) { 
            var $searchBar = $("#navbar-search-bar"); 
            if(!$searchBar.is(e.target) &&  
                $searchBar.has(e.target).length === 0) { 
                    $searchBar.removeClass("open");
            } 
        }); 

        // Click on the search button on navbar
        $("#btn-navbar-search").click(function(){
            var value = $navbarSearchInput.value;

            if(value.trim().length > 0) {
                var url = "/dashboard/company/users/search/all?query=" + encodeURI(value.trim());
                window.location = url;
            }
        });

        // Execute a function when the user releases a key on the keyboard
        $navbarSearchInput.addEventListener("keyup", function(event) {
            // Number 13 is the "Enter" key on the keyboard
            if (event.keyCode === 13) {
            // Cancel the default action, if needed
            event.preventDefault();
            // Trigger the button element with a click
            $("#btn-navbar-search").click();
            }
        }); 
    }
});

// Show/hide the loading screen
function toggleLoading(){    
    document.activeElement.blur();
    if (loadingOverlay.classList.contains('hidden')){
        loadingOverlay.classList.remove('hidden');
    } else {
        loadingOverlay.classList.add('hidden');
    }
}

// Get the search parameters in URL
function getSearchParameters() {
    var prmstr = window.location.search.substr(1);
    return prmstr != null && prmstr != "" ? transformToAssocArray(prmstr) : {};
}

function transformToAssocArray( prmstr ) {
  var params = {};
  var prmarr = prmstr.split("&");
  for ( var i = 0; i < prmarr.length; i++) {
      var tmparr = prmarr[i].split("=");
      params[tmparr[0]] = tmparr[1];
  }
  return params;
}
