function drawChart(inp) {
    var options = {
        title: 'Viewership',
        legend: { position: 'bottom' },
        enableInteractivity: 'true',
        explorer: {axis:'horizontal'}
    };
    var data = google.visualization.arrayToDataTable(inp)

    var chart = new google.charts.Line(document.getElementById('curve_chart'));
    chart.draw(data, google.charts.Line.convertOptions(options));
}

function listFollowing(data) {
        for(i = 0 ; i < data.length;i++) {
            $('#list_streams').append('<option>' + data[i] +  '</option>');
        }
        $('#loading').hide();
        $('#div_list_streams').show();
}
