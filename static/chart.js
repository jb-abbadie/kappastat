function drawChart(inp) {
    var options = {
        title: 'Viewership',
        curveType: 'none',
        legend: { position: 'bottom' },
        enableInteractivity: 'true',
        explorer: {axis:'horizontal'}
    };
    var data = google.visualization.arrayToDataTable(inp)

    var chart = new google.visualization.LineChart(document.getElementById('curve_chart'));
    chart.draw(data, options);
}


$(document).ready(function() {
    $('#curve_chart').height('80%');
    $('#curve_chart').width('100%');
    $('#list_streams').hide();
    $.getJSON("/api/following", function( data) {

        for(i = 0 ; i < data.length;i++) {
            $('#list_streams').append('<option value="' + data[i] + '">' + data[i] + '</div>');
        }
        $('#loading').hide();
        $('#list_streams').show();
    });


    $('#list_streams').change( function() {
        $('#curve_chart').fadeOut(500);
        $.getJSON("/api/viewer/" + $('#list_streams option:selected').text(), function(data) {
            console.log(data[0]);

            var result = [];
            result.push(['Time', 'Viewer']);
            for(var i in data) {
                console.log(data[i]);
                result.push([new Date(data[i]['Time']), data[i]['Viewer']]);
            }
            $('#curve_chart').fadeIn(1500);
            drawChart(result);
            console.log(result);
        });
    });
});

