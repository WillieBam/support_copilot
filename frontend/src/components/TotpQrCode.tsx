import { useEffect, useState } from 'react'
import QRCode from 'qrcode'

type TotpQrCodeProps = {
  otpauthUrl: string
}

function TotpQrCode({ otpauthUrl }: TotpQrCodeProps) {
  const [dataUrl, setDataUrl] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true

    const renderQr = async () => {
      if (!otpauthUrl) {
        setDataUrl('')
        setError('')
        return
      }

      try {
        const url = await QRCode.toDataURL(otpauthUrl, {
          width: 220,
          margin: 1,
        })
        if (!isMounted) return
        setDataUrl(url)
        setError('')
      } catch {
        if (!isMounted) return
        setDataUrl('')
        setError('Unable to generate QR image. You can still copy the setup URL below.')
      }
    }

    void renderQr()

    return () => {
      isMounted = false
    }
  }, [otpauthUrl])

  if (!otpauthUrl) return null

  return (
    <div className="totp-qr-wrap">
      {dataUrl ? (
        <img className="totp-qr-image" src={dataUrl} alt="TOTP setup QR code" />
      ) : (
        <div className="totp-qr-placeholder">Preparing QR image...</div>
      )}
      {error ? <p className="auth-help">{error}</p> : null}
    </div>
  )
}

export default TotpQrCode
