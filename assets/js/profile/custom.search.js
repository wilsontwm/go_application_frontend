$(document).ready(function(){   
    function profileTemplate(data) {
        var html = '';
        $.each(data, function(index, item){
            html += '<tr>'
                    + '<td>'
                    + '<a href="/dashboard/user/' + item["ID"] + '" class="dropdown-item">'
                    + '<div class="row">'
                    + '<img src="' + item["profilePicture"] + '" alt="' + item["name"] + '" class="img-size-50 mr-3 img-circle" style="height:50px;">'
                    + '<div class="">'
                    + '<p class="text-sm">' + item["name"] + '</p>'
                    + '<p class="text-sm text-muted">' + item["email"] + '</p>'
                    + '</div>'
                    + '</div>'
                    + '</a>'
                    + '</td>'
                    + '</tr>';
        });

        return html;
    }
    
    // Get the profile search result via AJAX 
    function loadProfileAjax(url) {
        axios.get(url)
        .then(function (response) {
            // handle success        
            data = response["data"]["data"];

            if(data.length > 0) {
                $('#profile-search-pagination-container').pagination({
                    pageSize: 25,
                    showGoInput: true,
                    showGoButton: true,
                    dataSource: data,
                    callback: function(d, pagination) {
                        // template method of yourself
                        var html = profileTemplate(d);
                        $('#profile-search-results-container').html(html);
                    }
                }); 
            } else {
                $('#profile-search-results').html('<h4 class="text-center text-muted">Ops.. No results found.</h4>');
            }
             
        })
        .catch(function (error) {
            // handle error
            console.log(error);
        });    
    }

    var $searchpane = document.getElementById("profile-search-pane");
    if (typeof($searchpane) != 'undefined' && $searchpane != null) {
        // Get the URL from the request
        var params = getSearchParameters();
        var query = decodeURI(params.query);
        $('#input-navbar-search').val(query);
        $('#search-query').html('"' + query + '"');

        var url = "/dashboard/company/users/search" + window.location.search;
        loadProfileAjax(url);
    }
});