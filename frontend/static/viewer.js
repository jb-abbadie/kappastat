$(document).ready(function() {
    $('#curve_chart').height('80%');
    $('#curve_chart').width('100%');
    $('#list_streams').hide();
    $.getJSON("/api/following", listFollowing);

    $('#list_streams').change( function() {
        $('#curve_chart').fadeOut(500, function() { drawChart(result);$('#curve_chart').fadeIn(1500);});
        $.getJSON("/api/viewer/" + $('#list_streams option:selected').text(), function(data) {

            var result = [];
            result.push(['Time', 'Viewer']);
            for(var i in data) {
                result.push([new Date(data[i]['Time']), data[i]['Viewer']]);
            }
        });
    });
});

