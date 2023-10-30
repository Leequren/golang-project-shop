create table display (
    idDisplay serial not null primary key,
    NameDisplay varchar,
    LengthDiagonalDisplay decimal(4, 1) not null,
    ResolutionDisplay varchar not null,
    MatrixDisplay varchar not null,
    UseGSync bool not null
);

create table monitor (
    idMonitor serial not null primary key,
    NameMonitor varchar not null,
    VoltageMonitor decimal(4, 1) not null,
    UseGSyncPremMonitor bool not null,
    CurvedMonitor bool not null,
    DisplayMonitorId int not null references display(idDisplay)
)