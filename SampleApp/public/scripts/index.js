
const addData = async (event) => {
    event.preventDefault();

    const carId = document.getElementById("carId").value;
    const make = document.getElementById("make").value;
    const model = document.getElementById("model").value;
    const color = document.getElementById("color").value;
    const dateOfManufacture =
        document.getElementById("dateOfManufacture").value;
    const manufacturerName =
        document.getElementById("manufacturerName").value;

    const carData = {
        carId: carId,
        make: make,
        model: model,
        color: color,
        dateOfManufacture: dateOfManufacture,
        manufacturerName: manufacturerName,
    };

    if (
        carId.length === 0 ||
        make.length === 0 ||
        model.length === 0 ||
        color.length === 0 ||
        dateOfManufacture.length === 0 ||
        manufacturerName.length === 0
    ) {
        alert("Please enter the data properly.");
        return;
    }

    try {
        const response = await fetch("/api/car", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(carData),
        });

        console.log("RESPONSE:", response);

        const data = await response.json();
        console.log("DATA:", data);

        if (!response.ok) {
            throw new Error(data.message || "Failed to create car.");
        }

        alert("Car Created");
    } catch (err) {
        alert("Error creating car.");
        console.error(err);
    }
};

const readData = async (event) => {
    event.preventDefault();

    const carId = document.getElementById("carIdInput").value.trim();

    if (carId.length === 0) {
        alert("Please enter a valid ID.");
        return;
    }

    try {
        const response = await fetch(
            `/api/car/${encodeURIComponent(carId)}`
        );

        const responseData = await response.json();

        console.log("Response data:", responseData);

        if (!response.ok) {
            throw new Error(
                responseData.message || "Failed to retrieve car."
            );
        }

        alert(JSON.stringify(responseData, null, 2));
    } catch (err) {
        alert("Error retrieving car.");
        console.error(err);
    }
};