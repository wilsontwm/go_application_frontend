$(document).ready(function(){    
    $invitationForm = document.getElementById("respond-company-invitation-form");
    $(".btn-view-invitation").click(function(){
        var id = $(this).data('id');
        var companyName = $(this).data('company');
        var status = $(this).data('status');
        $invitationForm.action = "/dashboard/invite/incoming/" + id + "/respond";
        $("#company-name-container").html(companyName);
        if(parseInt(status) == 0) {
            $("#btn-accept-invitation").attr("disabled", false);
            $("#btn-decline-invitation").attr("disabled", false);
        } else {
            $("#btn-accept-invitation").attr("disabled", true);
            $("#btn-decline-invitation").attr("disabled", true);
        }
        $("#modal-view-invitation").modal("show");
    });
    
    $("#btn-accept-invitation").click(function(e) {
        e.preventDefault();
        $("#inputResponse").val("accept");
        $invitationForm.submit();
        $("#modal-view-invitation").modal("hide");
        toggleLoading();
    });

    $("#btn-decline-invitation").click(function(e) {
        e.preventDefault();
        $("#inputResponse").val("decline");
        $invitationForm.submit();
        $("#modal-view-invitation").modal("hide");
        toggleLoading();
    });

});
