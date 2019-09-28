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
