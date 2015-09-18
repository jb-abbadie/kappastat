$(function() {
    $('#curve_chart').height('80%');
    $('#curve_chart').width('100%');
    $('#list_streams').hide();
    $.getJSON("/api/following", listFollowing);

    $('#list_streams').change( function() {
        $('#curve_chart').fadeOut(500);
        $.getJSON("/api/stat/" + $('#list_streams option:selected').text(), function(data) {

            var result = [];
            for(var i in data) {
                result.push([new Date(data[i]['Start']), data[i]['Viewer'], data[i]['Messages']]);
            }
            $('#curve_chart').fadeIn(1500);
            drawChart2(result);
        });
    });
});

function drawChart2(inp) {
    var options = {
        chart: {
            title: 'Test',
        },
        series : {
            0:{axis: 'Viewer'},
            1:{axis: 'Messages'}
        },
        axes: {
            y: {
                Viewer: {label: "Viewership"},
                Messages: {label: "Message"}
            }
        },
    };
    var data = new google.visualization.DataTable();
    data.addColumn('date', 'Time');
    data.addColumn('number', "Viewer");
    data.addColumn('number', "Chat Messages");
    data.addRows(inp);

    var chart = new google.charts.Line(document.getElementById('curve_chart'));
    chart.draw(data, google.charts.Line.convertOptions(options));
}
