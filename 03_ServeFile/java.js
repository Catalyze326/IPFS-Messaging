$(document).ready(function() {
    $("#lesen").click(function() {
        $.ajax({
            url : "mess.txt",
            dataType: "text",
            success : function (data) {
                $(".text").html(data);
            }
        });
    });
});
