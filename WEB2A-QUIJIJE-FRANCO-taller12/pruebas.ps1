[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$B = if ($args.Length -gt 0) { $args[0] } else { "http://localhost:8080" }

function c ($descripcion, $esperado, $metodo, $ruta, $cuerpo) {
    Write-Host "`n### $descripcion [esperado: $esperado]" -ForegroundColor Cyan
    
    $params = @{
        Uri = "$B$ruta"
        Method = $metodo
        UseBasicParsing = $true
    }
    
    if ($null -ne $cuerpo) {
        $params.Body = $cuerpo
        $params.ContentType = "application/json"
    }

    try {
        $res = Invoke-WebRequest @params
        Write-Host "HTTP/1.1 $($res.StatusCode) $($res.StatusDescription)" -ForegroundColor Green
        Write-Host $res.Content
    } catch {
        $res = $_.Exception.Response
        if ($res) {
            Write-Host "HTTP/1.1 $([int]$res.StatusCode) $($res.StatusDescription)" -ForegroundColor Yellow
            $stream = $res.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($stream)
            Write-Host $reader.ReadToEnd()
        } else {
            Write-Host "Error de conexion: $_" -ForegroundColor Red
        }
    }
}

c "POST válido" "201" "POST" "/tickets" '{"titulo": "Impresora sin tóner", "prioridad":"alta"}'
c "GET lista" "200" "GET" "/tickets" $null
c "GET por id" "200" "GET" "/tickets/1" $null
c "GET id inexistente" "404" "GET" "/tickets/99" $null
c "GET id no numérico" "400" "GET" "/tickets/abc" $null
c "POST JSON roto" "400" "POST" "/tickets" '{"titulo": "x'
c "POST título corto" "400" "POST" "/tickets" '{"titulo": "abc", "prioridad":"alta"}'
c "POST prioridad mala" "400" "POST" "/tickets" '{"titulo": "Teclado dañado", "prioridad":"urgente"}'
c "DELETE no soportado" "405" "DELETE" "/tickets/1" $null
c "ruta inexistente" "404" "GET" "/usuarios" $null
c "pánico controlado" "500" "GET" "/explotar" $null
c "campo desconocido" "400" "POST" "/tickets" '{"titulo": "Pantalla rota", "prioridad":"alta", "estado": "cerrado"}'
c "el servidor sigue vivo" "200" "GET" "/tickets" $null
