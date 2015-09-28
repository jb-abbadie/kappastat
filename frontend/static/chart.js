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
            console.log(data[i])
            $('#list_streams').append('<option value="'+  data[i]['name']+'">' + data[i]['display_name'] +  '</option>');
        }
        $('#loading').hide();
        $('#list_streams').fadeIn({'duration':500});
        $('#list_duration').fadeIn({'duration':500});
        $('#go_button').fadeIn({'duration':500});
        $('#curve_chart').fadeIn({'duration':500});
}
