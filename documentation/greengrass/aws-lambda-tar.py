# Ref. https://github.com/uhppoted/uhppoted/discussions/17

import tarfile
import boto3
import io
import time

def make_tarfile(tarobj, files):
    def make_tarinfo(filename, filebody):
        file_tarinfo = tarfile.TarInfo(filename)
        file_tarinfo.size = len(filebody.encode('utf-8'))
        file_tarinfo.mtime = time.time()
        return(file_tarinfo)

    def make_fileobj(fileobj, filename):
        fileobj.name = filename
        fileobj.seek(0)
        return(fileobj)

    with tarfile.open(fileobj=tarobj, mode="w:gz") as tar:
        for (fname, body) in files:
            file_tarinfo = make_tarinfo(fname, body)
            file_obj = make_fileobj(io.BytesIO(body.encode('utf-8')), fname)
            logger.debug("Adding {}, size {} to tarfile".format(fname, file_tarinfo.size))
            tar.addfile(file_tarinfo, fileobj=file_obj)
    tarobj.seek(0)
    return(tarobj)